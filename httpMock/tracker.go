package httpMock

import (
	"sync"
	"testing"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/stretchr/testify/require"
)

type tracker struct {
	exp callSeq
	act callSeq
	// We have to use mutex in order to sync `act` across different possible go-routines that might affect it (in concurrent tests)
	mu sync.RWMutex

	t *testing.T
}

func newTracker(t *testing.T) *tracker {
	return &tracker{
		exp: callSeq{},
		act: callSeq{},
		t:   t,
	}
}

func (t *tracker) actCalls() callSeq {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.act
}

func (t *tracker) expect(call httpCall) {
	t.exp[call.uniqueId()] = call
}

func (t *tracker) registerCall(ctx context.Context) {
	call := httpCall{
		url: ctx.Request().URL.Path,
	}

	t.register(call)

	if _, ok := t.exp[call.uniqueId()]; !ok {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	call = t.exp[call.uniqueId()]

	ctx.StatusCode(call.respStatus)
	_, err := ctx.Write(call.respBody)
	require.NoError(t.t, err)
}

func (t *tracker) register(call httpCall) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if c, ok := t.act[call.uniqueId()]; ok {
		call.times = c.times
	}
	call.times += 1
	t.act[call.uniqueId()] = call
}
