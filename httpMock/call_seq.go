package httpMock

type callSeq map[uniqueId]httpCall

func (seq callSeq) stats() stats {
	result := stats{}
	for k, v := range seq {
		result[k] = v.Times
	}
	return result
}

type stats map[uniqueId]int
