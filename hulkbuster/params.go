package hulkbuster

import "github.com/tuneinsight/lattigo/v5/he/hefloat"

func GenParam(LogN int, depth int) hefloat.Parameters {
	var err error
	var params hefloat.Parameters

	INITIAL_PRIME := 55
	SCALE_FACTOR := 45
	KEY_SWITCHING_PRIME := 61
	logQ_list := []int{INITIAL_PRIME}
	for i := 0; i < depth; i++ {
		logQ_list = append(logQ_list, SCALE_FACTOR)
	}

	if params, err = hefloat.NewParametersFromLiteral(
		hefloat.ParametersLiteral{
			LogN:            LogN,                       // A ring degree of 2^{14}
			LogQ:            logQ_list,                  // An initial prime of 55 bits and 7 primes of 45 bits
			LogP:            []int{KEY_SWITCHING_PRIME}, // The log2 size of the key-switching prime
			LogDefaultScale: SCALE_FACTOR,               // The default log2 of the scaling factor
		}); err != nil {
		panic(err)
	}

	return params
}
