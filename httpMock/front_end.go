package httpMock

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kataras/iris"
	"github.com/stretchr/testify/require"
)

type FrontEnd struct {
	server *httptest.Server

	tracker *tracker

	exp callSeq

	t *testing.T
}

func New(t *testing.T) *FrontEnd {
	tracker := newTracker(t)

	i := iris.Default()
	i.OnErrorCode(iris.StatusNotFound, tracker.registerCall)

	require.NoError(t, i.Build())

	return &FrontEnd{
		server:  httptest.NewServer(i.Router),
		tracker: tracker,
		exp:     callSeq{},
		t:       t,
	}
}

func (c *FrontEnd) Mock(url string) *httpCallBuilder {
	return &httpCallBuilder{
		c:   c,
		url: url,
	}
}

func (c *FrontEnd) register(call httpCall) {
	id := call.uniqueId()
	if _, ok := c.exp[id]; ok {
		panic(fmt.Sprintf("call with id = `%v` was already registered", id))
	}
	c.reg(call)
}

func (c *FrontEnd) reRegister(call httpCall) {
	id := call.uniqueId()
	if _, ok := c.exp[id]; !ok {
		panic(fmt.Sprintf("missing Return() for call with id = `%v`", id))
	}
	c.reg(call)
}

func (c *FrontEnd) reg(call httpCall) {
	c.exp[call.uniqueId()] = call
	c.tracker.expect(call)
}

func (c *FrontEnd) Finish() {
	c.server.Close()

	c.compare(c.exp, c.tracker.actCalls())
}

func (c *FrontEnd) compare(expSeq, actSeq callSeq) (diff string, match bool) {
	exp := expSeq.stats()
	act := actSeq.stats()

	if !cmp.Equal(exp, act) {
		diff := fmt.Sprintf("expected:\n\t%v\nactual:\n\t%v\n", exp, act)
		require.Fail(c.t, diff)
	}

	return
}

func (c *FrontEnd) Host() string {
	return c.server.URL
}
