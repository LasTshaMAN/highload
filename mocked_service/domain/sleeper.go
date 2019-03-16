//go:generate mockgen -source sleeper.go -destination sleeper.mock.gen.go -package domain

package domain

import (
	"time"
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
	span, err := RandPointBetween(t1, t2)
	if err != nil {
		return err
	}
	time.Sleep(span)
	return nil
}

func (v *sleeperImpl) LoopForever() {
	for {
		time.Sleep(time.Millisecond)
	}
}
