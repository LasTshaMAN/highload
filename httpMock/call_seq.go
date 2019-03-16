package httpMock

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type callSeq map[uniqueId]httpCall

func (seq callSeq) stats() stats {
	result := stats{}
	for k, v := range seq {
		result[k] = v.times
	}
	return result
}

type stats map[uniqueId]int

func compare(expSeq, actSeq callSeq) (diff string, match bool) {
	match = true

	exp := expSeq.stats()
	act := actSeq.stats()

	if !cmp.Equal(exp, act) {
		//if !cmp.Equal(exp, act, cmp.AllowUnexported(uniqueId{})) {
		diff = fmt.Sprintf("expected:\n\t%v\nactual:\n\t%v\n", exp, act)
		match = false
	}

	return
}
