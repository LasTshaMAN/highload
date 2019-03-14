//go:generate mockgen -source sleeper.go -destination sleeper.mock.gen.go -package domain

package domain

import (
	"math/rand"
	"time"

	"github.com/kataras/iris/core/errors"
)

type Sleeper interface {
	Sleep(duration time.Duration)
	SleepInterval(t1 time.Duration, t2 time.Duration) error
	LoopForever()
}

func NewSleeper() *sleeperImpl {
	return &sleeperImpl{}
}

type sleeperImpl struct{}

func (v *sleeperImpl) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (v *sleeperImpl) SleepInterval(t1, t2 time.Duration) error {
	if t1 > t2 {
		return errors.New("t1 <= t2 must hold")
	}
	delta := t2 - t1
	randomSpan := rand.Intn(int(delta))
	time.Sleep(time.Duration(100+randomSpan) * time.Millisecond)
	return nil
}

func (v *sleeperImpl) LoopForever() {
	for {
		time.Sleep(time.Millisecond)
	}
}
