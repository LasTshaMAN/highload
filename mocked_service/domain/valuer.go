package domain

import "math/rand"

type Valuer interface {
	Value() int
}

func NewValuer() *valuerImpl {
	return &valuerImpl{}
}

type valuerImpl struct{}

func (v *valuerImpl) Value() int {
	return 1 + rand.Intn(100)
}
