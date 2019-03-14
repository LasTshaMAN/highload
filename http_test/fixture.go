package httpTest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type Fixture struct {
	t       *testing.T
	server  *httptest.Server
	toClose []io.ReadCloser
}

func NewFixture(t *testing.T, handler http.Handler) *Fixture {
	s := httptest.NewServer(handler)
	return &Fixture{
		t:      t,
		server: s,
	}
}

func (fx *Fixture) Finish() {
	for _, resp := range fx.toClose {
		require.NoError(fx.t, resp.Close())
	}
	fx.server.Close()
}

func (fx *Fixture) Request(method, url string, body io.Reader) *http.Response {
	r, err := http.NewRequest(method, fx.server.URL+url, body)
	if err != nil {
		require.NoError(fx.t, err)
	}
	resp, err := fx.server.Client().Do(r)
	require.NoError(fx.t, err)
	fx.toClose = append(fx.toClose, resp.Body)

	return resp
}

func (fx *Fixture) Parse(respBody io.Reader, result interface{}) {
	require.NoError(fx.t, json.NewDecoder(respBody).Decode(&result))
}
