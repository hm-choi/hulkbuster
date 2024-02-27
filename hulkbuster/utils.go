package hulkbuster

import (
	"fmt"
	"sync"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func RoutineRotation(rot int, input_ctxt *rlwe.Ciphertext, output_cipher []*rlwe.Ciphertext, wg *sync.WaitGroup, i int, eval hefloat.Evaluator) {
	eval.Rotate(input_ctxt, rot, output_cipher[i])
	defer wg.Done()
}

func MultAndRot(rot int, input_ctxt *rlwe.Ciphertext, coeff []float64, output_cipher []*rlwe.Ciphertext, wg *sync.WaitGroup, i int, eval hefloat.Evaluator) {
	ii := input_ctxt.CopyNew()
	if i != 0 {
		ii, _ = eval.RotateNew(input_ctxt, rot)
	}
	output_cipher[i], _ = eval.MulRelinNew(ii, coeff)
	eval.Rescale(output_cipher[i], output_cipher[i])
	defer wg.Done()
}

func RotAndMult(rot int, input_ctxt *rlwe.Ciphertext, coeff []float64, output_cipher []*rlwe.Ciphertext, wg *sync.WaitGroup, i int, eval hefloat.Evaluator) {
	ii, _ := eval.MulRelinNew(input_ctxt, coeff)
	if i != 0 {
		ii, _ = eval.RotateNew(ii, rot)
	}
	output_cipher[i] = ii
	eval.Rescale(output_cipher[i], output_cipher[i])
	defer wg.Done()
}

func RotateAndSum(rot int, input_ctxt *rlwe.Ciphertext, output_cipher []*rlwe.Ciphertext, wg *sync.WaitGroup, i int, eval hefloat.Evaluator) {
	eval.Rotate(input_ctxt, rot, output_cipher[i])
	defer wg.Done()
}
func RotateAndSum2(rot int, input_ctxt *rlwe.Ciphertext, output_cipher []*rlwe.Ciphertext, i int, eval hefloat.Evaluator) {
	eval.Rotate(input_ctxt, rot, output_cipher[i])
}

func KeyGen(params hefloat.Parameters, gls_list []int) (*rlwe.SecretKey, *rlwe.PublicKey, *rlwe.RelinearizationKey, *rlwe.MemEvaluationKeySet, []*rlwe.GaloisKey) {
	kgen := rlwe.NewKeyGenerator(params)
	sk := kgen.GenSecretKeyNew()
	pk := kgen.GenPublicKeyNew(sk)
	rlk := kgen.GenRelinearizationKeyNew(sk)
	evk := rlwe.NewMemEvaluationKeySet(rlk)
	galEls := []uint64{params.GaloisElementForComplexConjugation()}

	for idx := range gls_list {
		galEls = append(galEls, params.GaloisElement(gls_list[idx]))
	}
	glk := kgen.GenGaloisKeysNew(galEls, sk)
	return sk, pk, rlk, evk, glk
}

func OperatorGen(params hefloat.Parameters, pk *rlwe.PublicKey, rlk *rlwe.RelinearizationKey, evk *rlwe.MemEvaluationKeySet, glk []*rlwe.GaloisKey) (*hefloat.Encoder, *rlwe.Encryptor, *hefloat.Evaluator) {
	ecd := hefloat.NewEncoder(params)
	enc := rlwe.NewEncryptor(params, pk)
	eval := hefloat.NewEvaluator(params, evk).WithKey(rlwe.NewMemEvaluationKeySet(rlk, glk...))
	return ecd, enc, eval
}

func PrivOperatorGen(params hefloat.Parameters, sk *rlwe.SecretKey) *rlwe.Decryptor {
	return rlwe.NewDecryptor(params, sk)
}

func MultiplyMatrix(a [][]float64, b [][]float64) ([][]float64, error) {
	if len(a) == 0 || len(a[0]) == 0 || len(b) == 0 || len(b[0]) == 0 {
		return nil, fmt.Errorf("invalid array input")
	}
	if len(a[0]) != len(b) {
		return nil, fmt.Errorf("row length of first elem must be same as column length of second elem")
	}
	newMatrix := make([][]float64, len(a))
	for yIndex := range newMatrix {
		newMatrix[yIndex] = make([]float64, len(b[0]))
		for xIndex := range newMatrix[yIndex] {
			for iter := range b {
				newMatrix[yIndex][xIndex] += a[yIndex][iter] * b[iter][xIndex]
			}
		}
	}
	return newMatrix, nil
}

func MultiplyMatrix2(a [][]float64, b [][]float64) ([][]float64, error) {
	result_matrix := make([][]float64, 1)
	tmp := make([]float64, len(b))
	for i := 0; i < len(b); i++ {
		tmpp := 0.0
		for j := 0; j < len(a[0]); j++ {
			tmpp += a[0][j] * b[i][j]
		}
		tmp[i] = tmpp
	}
	result_matrix[0] = tmp
	return result_matrix, nil
}

func MatrixTransposeForFCLayer(a [][]float64, num_of_slots int) [][]float64 {
	result_vec := make([][]float64, len(a))
	for i := 0; i < len(result_vec); i++ {
		result_vec[i] = make([]float64, len(a)+len(a[0]))
	}
	fmt.Println("len(a), len(result_vec), len(result_vec[0])", len(a), len(result_vec), len(result_vec[0]))
	extend_vec := make([][]float64, len(a))
	for i := 0; i < len(a); i++ {
		extend_vec[i] = make([]float64, len(a)-i)
		extend_vec[i] = append(extend_vec[i], a[i]...)
		// result_vec[i] = append(result_vec[i], make([]float64, i)...)
	}
	fmt.Println("len(extend_vec), len(extend_vec[0])", len(extend_vec), len(extend_vec[0]))

	for i := 0; i < len(extend_vec); i++ {
		for j := 0; j < len(extend_vec[i]); j++ {
			h_idx := j % len(extend_vec)
			v_idx := j / len(extend_vec)
			result_vec[h_idx][v_idx] = extend_vec[i][j]
		}
	}

	for i := 0; i < len(extend_vec); i++ {
		tmp := result_vec[i][len(a):]
		tmp = append(tmp, make([]float64, (num_of_slots-(len(a)+len(tmp))))...)
		tmp = append(tmp, result_vec[i][:len(a)]...)
		result_vec[i] = tmp
	}

	return result_vec
}

func EvalListGen(eval *hefloat.Evaluator, num_thread int) []*hefloat.Evaluator {
	eval_list := make([]*hefloat.Evaluator, num_thread)
	eval_list[0] = eval
	if num_thread > 1 {
		for i := 1; i < num_thread; i++ {
			eval_list[i] = eval.ShallowCopy()
		}
	}
	return eval_list
}
