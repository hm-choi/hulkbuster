// package main

// import (
// 	"encoding/csv"
// 	"fmt"
// 	"math"
// 	"math/rand"
// 	"os"
// 	"strconv"
// 	"sync"
// 	"time"

// 	"github.com/hm-choi/hulkbuster"
// 	"github.com/tuneinsight/lattigo/v5/he/hefloat"
// )

// func main() {
// 	params := hulkbuster.GenParam(14, 5)
// 	fmt.Println("num of slots", params.MaxSlots())
// 	rot_list := []int{-10, -20, -30, -40, -50, -120}
// 	for i := 0; i < 120; i++ {
// 		rot_list = append(rot_list, i)
// 	}
// 	for i := 0; i < 17; i++ {
// 		rot_list = append(rot_list, i*120)
// 	}
// 	for i := 1; i < 8; i++ {
// 		rot_list = append(rot_list, i*(-1200))
// 	}
// 	sk, pk, rlk, evk, glk := hulkbuster.KeyGen(params, rot_list)
// 	dec := hulkbuster.PrivOperatorGen(params, sk)
// 	ecd, enc, eval := hulkbuster.OperatorGen(params, pk, rlk, evk, glk)
// 	NUM_THREAD := 8
// 	pt := hefloat.NewPlaintext(params, params.MaxLevel())
// 	eval_list := hulkbuster.EvalListGen(eval, NUM_THREAD)

// 	var wg sync.WaitGroup

// 	a := make([][]float64, 1)
// 	tmp := make([]float64, 1014)
// 	for j := 0; j < 1014; j++ {
// 		tmp[j] = rand.Float64()
// 		tmp[j] = 0.1
// 	}
// 	a[0] = tmp
// 	b := [][]float64{}
// 	for i := 0; i < 120; i++ {
// 		tmp := make([]float64, 1014)
// 		b = append(b, tmp)
// 	}

// 	fd, _ := os.Open("wow.csv")
// 	fileReader := csv.NewReader(fd)
// 	c, _ := fileReader.ReadAll()
// 	for i := 0; i < 120; i++ {
// 		for j := 0; j < 1014; j++ {
// 			b[i][j], _ = strconv.ParseFloat(c[i][j], 64)
// 		}
// 	}

// 	bb := make([][]float64, 120)
// 	for i := 0; i < 120; i++ {
// 		bb[i] = make([]float64, 8192)
// 	}
// 	fd, _ = os.Open("wow2.csv")
// 	fileReader = csv.NewReader(fd)
// 	c, _ = fileReader.ReadAll()
// 	for i := 0; i < 120; i++ {
// 		for j := 0; j < 8192; j++ {
// 			bb[i][j], _ = strconv.ParseFloat(c[i][j], 64)
// 		}
// 	}
// 	res_mat, _ := hulkbuster.MultiplyMatrix2(a, b)

// 	pt = hefloat.NewPlaintext(params, params.MaxLevel())

// 	ecd.Encode(a[0], pt)
// 	input_ctxt, _ := enc.EncryptNew(pt)

// 	START_TIME := time.Now()

// 	result := hulkbuster.FC_Layer(eval_list, &wg, input_ctxt, bb, NUM_THREAD, 1014, 120)
// 	fmt.Println("Time", time.Since(START_TIME))

// 	values2 := make([]float64, 1<<params.LogMaxSlots())
// 	ecd.Decode(dec.DecryptNew(result), values2)

// 	Sum := 0.0
// 	for i := 0; i < 120; i++ {
// 		Sum += math.Abs(res_mat[0][i] - values2[i])
// 	}
// 	fmt.Println("오차", Sum)
// }
