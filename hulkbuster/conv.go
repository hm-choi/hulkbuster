package hulkbuster

import (
	"sync"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func Conv_Layer(eval_list []*hefloat.Evaluator, wg *sync.WaitGroup, input_ctxt *rlwe.Ciphertext,
	coeff [][]float64, num_of_gorutine int) []*rlwe.Ciphertext {

	return nil
}
