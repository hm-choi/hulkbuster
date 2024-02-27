package hulkbuster

import (
	"sync"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func FC_Layer(eval_list []*hefloat.Evaluator, wg *sync.WaitGroup, input_ctxt *rlwe.Ciphertext,
	coeff [][]float64, num_thread int, input_size int, output_size int) *rlwe.Ciphertext {
	if len(coeff) < num_thread {
		num_thread = output_size
	}

	CON1, CON2, CON3 := 20, 10, 6
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
				go RotAndMult(IDX, input_ctxt, coeff[IDX], output_ciphers, wg, IDX, *eval_list[IDX%num_thread])
				IDX += 1
			}
			wg.Wait()
		} else {
			wg.Add(TMP)
			for j := CON1 - TMP; j < CON1; j++ {
				go RotAndMult(IDX, input_ctxt, coeff[IDX], output_ciphers, wg, IDX, *eval_list[IDX%num_thread])
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
		go RotateAndSum(i*1200, rotated_sum, output_ciphers, wg, i, *eval_list[i%num_thread])
	}
	wg.Wait()

	for i := 1; i < len(output_ciphers); i++ {
		eval_list[0].Add(output_ciphers[0], output_ciphers[i], output_ciphers[0])
	}

	return output_ciphers[0]
}
