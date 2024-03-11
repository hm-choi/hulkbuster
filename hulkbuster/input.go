package hulkbuster

import (
	"sync"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func InputCtxtExpansion(eval_list []*hefloat.Evaluator, wg *sync.WaitGroup, input_ctxt *rlwe.Ciphertext,
	ctxt_num int, rot int) *rlwe.Ciphertext {
	ctxt_list := make([]*rlwe.Ciphertext, ctxt_num)

	for i := 0; i < ctxt_num; i++ {
		ctxt_list[i] = input_ctxt.CopyNew()
	}

	wg.Add(ctxt_num)
	for i := 0; i < ctxt_num; i++ {
		rot_ := rot * (i + 1) * (-1)
		go RoutineRotation(rot_, input_ctxt, ctxt_list, wg, i, *eval_list[i])
	}
	wg.Wait()

	for i := 1; i < ctxt_num; i++ {
		eval_list[0].Add(input_ctxt, ctxt_list[i], input_ctxt)
	}

	return input_ctxt // ctxt_list[0]
}
