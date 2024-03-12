package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/hm-choi/hulkbuster"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func main() {
	params := hulkbuster.GenParam(14, 5)
	number_of_slots := params.MaxSlots()
	fmt.Println("num of slots", number_of_slots)

	input_size, output_size := 1014, 120
	data_size := output_size * int(math.Ceil(float64(input_size) / float64(output_size)) + 1)

	CON3 := number_of_slots / data_size // 6
	CON1 := output_size / CON3 			// 20
	CON2 := data_size / output_size 	// 10

	rot_list := []int{-1 * input_size}
	for i := 0; i < CON1; i++ {
		rot_list = append(rot_list, i)
	}
	for i := 1; i < CON2; i++ {
		rot_list = append(rot_list, i*output_size)
	}
	for i := 1; i < CON3; i++ {
		rot_list = append(rot_list, i * data_size)
	}
	sk, pk, rlk, evk, glk := hulkbuster.KeyGen(params, rot_list)
	dec := hulkbuster.PrivOperatorGen(params, sk)
	ecd, enc, eval := hulkbuster.OperatorGen(params, pk, rlk, evk, glk)

	NUM_THREAD := 8
	pt := hefloat.NewPlaintext(params, params.MaxLevel())
	eval_list := hulkbuster.EvalListGen(eval, NUM_THREAD)

	var wg sync.WaitGroup

	a := make([][]float64, 1)
	a[0] = make([]float64, input_size)
	for j := 0; j < input_size; j++ {
		a[0][j] = rand.Float64()
	}

	b := [][]float64{}
	for i := 0; i < output_size; i++ {
		tmp := make([]float64, input_size)
		b = append(b, tmp)
	}

	fd, _ := os.Open("wow3.csv")
	fileReader := csv.NewReader(fd)
	c, _ := fileReader.ReadAll()
	for i := 0; i < output_size; i++ {
		for j := 0; j < input_size; j++ {
			b[i][j], _ = strconv.ParseFloat(c[i][j], 64)
		}
	}

	bb := make([][]float64, CON1)
	for i := 0; i < CON1; i++ {
		bb[i] = make([]float64, number_of_slots)
	}
	fd, _ = os.Open("wow4.csv")
	fileReader = csv.NewReader(fd)
	c, _ = fileReader.ReadAll()
	for i := 0; i < CON1; i++ {
		for j := 0; j < data_size * CON3; j++ {
			bb[i][j], _ = strconv.ParseFloat(c[i][j], 64)
		}
	}
	res_mat, _ := hulkbuster.MultiplyMatrix2(a, b)

	a_plus := []float64{}

	for i := 0; i < CON3; i++ {
		a_plus = append(a_plus, make([]float64, output_size - i*CON1)...)
		a_plus = append(a_plus, a[0]...)
		a_plus = append(a_plus, make([]float64, data_size - input_size - output_size + i*CON1)...)
	}

	pt = hefloat.NewPlaintext(params, params.MaxLevel())

	ecd.Encode(a_plus, pt)
	input_ctxt, _ := enc.EncryptNew(pt)

	START_TIME := time.Now()

	result := hulkbuster.FC_Layer(eval_list, &wg, input_ctxt, bb, NUM_THREAD, input_size, output_size, number_of_slots)
	fmt.Println("Time", time.Since(START_TIME))

	values2 := make([]float64, 1<<params.LogMaxSlots())
	ecd.Decode(dec.DecryptNew(result), values2)

	Sum := 0.0
	for i := 0; i < output_size; i++ {
		Sum += math.Abs(res_mat[0][i] - values2[i])
		// if i < CON2 {
		// 	fmt.Println(res_mat[0][i], values2[i])
		// }
	}
	fmt.Println("오차", Sum)
}
