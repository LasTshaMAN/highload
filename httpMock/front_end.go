package httpMock

import (
	"fmt"
	"net/http/httptest"
	"testing"

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

	if diff, match := compare(c.exp, c.tracker.actCalls()); !match {
		require.Fail(c.t, diff)
	}
}

func (c *FrontEnd) Host() string {
	return c.server.URL
}
