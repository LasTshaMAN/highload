package httpMock

import (
	"testing"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/stretchr/testify/require"
)

type tracker struct {
	exp callSeq
	act callSeq

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
	return t.act
}

func (t *tracker) expect(call httpCall) {
	t.exp[call.uniqueId()] = call
}

func (t *tracker) registerCall(ctx context.Context) {
	call := httpCall{
		URL: ctx.Request().URL.Path,
	}

	t.register(call)

	if _, ok := t.exp[call.uniqueId()]; !ok {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	call = t.exp[call.uniqueId()]

	ctx.StatusCode(call.RespStatus)
	_, err := ctx.Write(call.RespBody)
	require.NoError(t.t, err)
}

func (t *tracker) register(call httpCall) {
	if c, ok := t.act[call.uniqueId()]; ok {
		call.Times = c.Times
	}
	call.Times += 1
	t.act[call.uniqueId()] = call
}
