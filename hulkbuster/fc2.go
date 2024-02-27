package hulkbuster

// import (
// 	"fmt"
// 	"sync"

// 	"github.com/tuneinsight/lattigo/v5/core/rlwe"
// 	"github.com/tuneinsight/lattigo/v5/he/hefloat"
// )

// func FC_Layer(eval_list []*hefloat.Evaluator, wg *sync.WaitGroup, input_ctxt *rlwe.Ciphertext,
// 	coeff [][]float64, num_thread int, input_size int, output_size int) *rlwe.Ciphertext {
// 	if len(coeff) < num_thread {
// 		num_thread = output_size
// 	}

// 	output_ciphers := make([]*rlwe.Ciphertext, output_size)
// 	for i := 0; i < output_size; i++ {
// 		output_ciphers[i] = input_ctxt.CopyNew()
// 	}
// 	for i := 0; i < output_size/num_thread; i++ {
// 		wg.Add(num_thread)
// 		for j := 0; j < num_thread; j++ {
// 			idx := j + i*num_thread
// 			go MultAndRot(idx, input_ctxt, coeff[idx], output_ciphers, wg, idx, *eval_list[j])
// 		}
// 		wg.Wait()
// 	}

// 	// tmp, _ := eval_list[0].MulRelinNew(input_ctxt, coeff[0])
// 	// return tmp
// 	// Gorutine이 더 빠르면 Gorutine 적용
// 	for i := 1; i < len(output_ciphers); i++ {
// 		eval_list[0].Add(output_ciphers[0], output_ciphers[i], output_ciphers[0])
// 	}

// 	rotated_sum := output_ciphers[0]
// 	// return rotated_sum
// 	output_s := input_size / output_size
// 	fmt.Println(output_s)
// 	output_ciphers = make([]*rlwe.Ciphertext, output_s+1)
// 	for i := 0; i < output_s+1; i++ {
// 		output_ciphers[i] = rotated_sum.CopyNew()
// 	}

// 	for idx := 0; idx < 9; idx++ {
// 		wg.Add(1)
// 		if idx == 0 {
// 			go RotateAndSum((-1)*output_size, rotated_sum, output_ciphers, wg, idx, *eval_list[idx%num_thread])
// 		} else {
// 			go RotateAndSum(idx*output_size, rotated_sum, output_ciphers, wg, idx, *eval_list[idx%num_thread])
// 		}
// 	}
// 	wg.Wait()

// 	// Gorutine이 더 빠르면 Gorutine 적용
// 	for i := 0; i < len(output_ciphers); i++ {
// 		eval_list[0].Add(rotated_sum, output_ciphers[i], rotated_sum)
// 	}

// 	return rotated_sum
// }
