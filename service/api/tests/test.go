package tests

import (
	"highload/service/api"
	"highload/service/api/responses"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type serverBootstrap func(t *testing.T, a *api.API) http.Handler
type testCase func(t *testing.T, bootstrap serverBootstrap)

func NewTestSuite() map[string]testCase {
	return map[string]testCase{
		"/api/endpoint": func(t *testing.T, sb serverBootstrap) {
			fx := newFixture(t, sb)
			defer fx.ctr.Finish()
			defer fx.Finish()

			expValue := 50
			fx.avg.EXPECT().Value().Return(expValue)

			resp := fx.Request("GET", "/api/endpoint", nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			var body responses.Answer
			fx.Parse(resp.Body, &body)
			require.Equal(t, expValue, body.Value)
		},
	}
}
