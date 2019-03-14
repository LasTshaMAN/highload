package tests

import (
	"fmt"
	"highload/service/api/responses"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase func(t *testing.T)

func NewTestSuite(handler http.Handler) map[string]testCase {
	return map[string]testCase{
		"/api/fast": func(t *testing.T) {
			fx := newFixture(t, handler)
			defer fx.finish()

			resp := fx.request("GET", "/api/fast", nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			var body responses.Answer
			fx.parse(resp.Body, &body)
			require.True(t, 0 < body.Value && body.Value < 101, fmt.Sprintf("value = %d", body.Value))
		},
	}
}
