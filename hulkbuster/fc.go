package hulkbuster

import (
	"sync"
	"math"
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func FC_Layer(eval_list []*hefloat.Evaluator, wg *sync.WaitGroup, input_ctxt *rlwe.Ciphertext,
	coeff [][]float64, num_thread int, input_size int, output_size int, number_of_slots int) *rlwe.Ciphertext {
	if len(coeff) < num_thread {
		num_thread = output_size
	}

	data_size := output_size * int(math.Ceil(float64(input_size) / float64(output_size)) + 1)
	CON3 := number_of_slots / data_size // 6
	CON1 := output_size / CON3 			// 20
	CON2 := data_size / output_size 	// 10

	output_ciphers := make([]*rlwe.Ciphertext, CON1)
	for i := 0; i < CON1; i++ {
		output_ciphers[i] = input_ctxt.CopyNew()
	}

	IDX := 0
	for IDX < CON1 {
		TMP := CON1 - IDX
		if TMP >= num_thread {
			wg.Add(num_thread)
			for j := 0; j < num_thread; j++ {
				go MultAndRot(IDX, input_ctxt, coeff[IDX], output_ciphers, wg, IDX, *eval_list[IDX%num_thread])
				IDX += 1
			}
			wg.Wait()
		} else {
			wg.Add(TMP)
			for j := CON1 - TMP; j < CON1; j++ {
				go MultAndRot(IDX, input_ctxt, coeff[IDX], output_ciphers, wg, IDX, *eval_list[IDX%num_thread])
				IDX += 1
			}
			wg.Wait()
		}
	}

	// Gorutine이 더 빠르면 Gorutine 적용
	for i := 1; i < len(output_ciphers); i++ {
		eval_list[0].Add(output_ciphers[0], output_ciphers[i], output_ciphers[0])
	}

	rotated_sum := output_ciphers[0].CopyNew()

	// output_s := input_size / output_size
	output_ciphers = make([]*rlwe.Ciphertext, CON2)
	for i := 0; i < CON2; i++ {
		output_ciphers[i] = rotated_sum.CopyNew()
	}

	IDX = 0
	for IDX < CON2 {
		TMP := CON2 - IDX
		if TMP >= num_thread {
			wg.Add(num_thread)
			for j := 0; j < num_thread; j++ {
				go RotateAndSum((IDX)*output_size, rotated_sum, output_ciphers, wg, IDX, *eval_list[IDX%num_thread])
				IDX += 1
			}
			wg.Wait()
		} else {
			wg.Add(TMP)
			for j := CON2 - TMP; j < CON2; j++ {
				go RotateAndSum((IDX)*output_size, rotated_sum, output_ciphers, wg, IDX, *eval_list[IDX%num_thread])
				IDX += 1
			}
			wg.Wait()
		}
	}

	// Gorutine이 더 빠르면 Gorutine 적용
	for i := 1; i < len(output_ciphers); i++ {
		eval_list[0].Add(output_ciphers[0], output_ciphers[i], output_ciphers[0])
	}

	rotated_sum = output_ciphers[0].CopyNew()

	output_ciphers = make([]*rlwe.Ciphertext, CON3)
	for i := 0; i < CON3; i++ {
		output_ciphers[i] = rotated_sum.CopyNew()
	}

	wg.Add(CON3)
	for i := 0; i < CON3; i++ {
		go RotateAndSum(i*data_size, rotated_sum, output_ciphers, wg, i, *eval_list[i%num_thread])
	}
	wg.Wait()

	for i := 1; i < len(output_ciphers); i++ {
		eval_list[0].Add(output_ciphers[0], output_ciphers[i], output_ciphers[0])
	}

	return output_ciphers[0]
}
