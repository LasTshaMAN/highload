package tests

import (
	"errors"
	"fmt"
	"highload/mocked_service/api"
	"highload/mocked_service/api/responses"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type serverBootstrap func(t *testing.T, a *api.API) http.Handler
type testCase func(t *testing.T, bootstrap serverBootstrap)

func NewTestSuite() map[string]testCase {
	return map[string]testCase{
		"/api/fast": func(t *testing.T, sb serverBootstrap) {
			fx := newFixture(t, sb)
			defer fx.ctr.Finish()
			defer fx.Finish()

			resp := fx.Request("GET", "/api/fast", nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			var body responses.Answer
			fx.Parse(resp.Body, &body)
			require.True(t, 0 < body.Value && body.Value < 101, fmt.Sprintf("value = %d", body.Value))
		},
		"/api/slow": func(t *testing.T, sb serverBootstrap) {
			fx := newFixture(t, sb)
			defer fx.ctr.Finish()
			defer fx.Finish()

			fx.sleeper.EXPECT().Sleep(1000 * time.Millisecond)

			resp := fx.Request("GET", "/api/slow", nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			var body responses.Answer
			fx.Parse(resp.Body, &body)
			require.True(t, 0 < body.Value && body.Value < 101, fmt.Sprintf("value = %d", body.Value))
		},
		"/api/rand": func(t *testing.T, sb serverBootstrap) {
			fx := newFixture(t, sb)
			defer fx.ctr.Finish()
			defer fx.Finish()

			fx.sleeper.EXPECT().SleepInterval(100*time.Millisecond, 1000*time.Millisecond)

			resp := fx.Request("GET", "/api/random", nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			var body responses.Answer
			fx.Parse(resp.Body, &body)
			require.True(t, 0 < body.Value && body.Value < 101, fmt.Sprintf("value = %d", body.Value))
		},
		"/api/rand (error case)": func(t *testing.T, sb serverBootstrap) {
			fx := newFixture(t, sb)
			defer fx.ctr.Finish()
			defer fx.Finish()

			fx.sleeper.EXPECT().SleepInterval(100*time.Millisecond, 1000*time.Millisecond).Return(errors.New(""))

			resp := fx.Request("GET", "/api/random", nil)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		},
		"/api/never_ending": func(t *testing.T, sb serverBootstrap) {
			fx := newFixture(t, sb)
			defer fx.ctr.Finish()
			defer fx.Finish()

			fx.sleeper.EXPECT().LoopForever()

			resp := fx.Request("GET", "/api/never_ending", nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		},
	}
}
