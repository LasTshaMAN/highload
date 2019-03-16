package domain

import (
	"errors"
	"math/rand"
	"time"
)

func RandPointBetween(t1, t2 time.Duration) (time.Duration, error) {
	if t1 > t2 {
		return 0, errors.New("t1 <= t2 must hold")
	}
	fraction := rand.Float64()
	return t1 + time.Duration(fraction*float64(t2-t1)), nil
}
