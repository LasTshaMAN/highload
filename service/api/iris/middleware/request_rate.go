package middleware

import (
	"log"
	"sync/atomic"
	"time"
)

func LogRequestRate() {
	count := uint64(0)
	tick := time.NewTicker(time.Second)
	atomic.AddUint64(&count, 1)
	select {
	case <-tick.C:
		log.Printf("Rate: %d / sec", atomic.SwapUint64(&count, 0))
	default:
	}
}
