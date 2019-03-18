package httpMock

import (
	"crypto/tls"
	"fmt"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/http2"

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

	// TODO
	// What do we do about mocking HTTP 1.x (maybe we need to create a separate mock for it ...) ?

	// we need to configure TLS on test server (otherwise it falls back to HTTP 1.1)
	// for more details see "http://big-elephants.com/2017-09/this-programmer-tried-to-mock-an-http-slash-2-server-in-go-and-heres-what-happened"
	h2Server := httptest.NewUnstartedServer(i.Router)
	h2Server.TLS = &tls.Config{
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		NextProtos:   []string{http2.NextProtoTLS},
	}
	h2Server.StartTLS()

	return &FrontEnd{
		server:  h2Server,
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
