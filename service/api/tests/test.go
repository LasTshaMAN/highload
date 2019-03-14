package tests

import (
	"fmt"
	"highload/http_test"
	"highload/service/api"
	"highload/service/api/responses"
	"highload/service/domain"
	"net/http"
	"testing"
	"time"

	"github.com/kataras/iris/core/errors"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"
)

type serverBootstrap func(t *testing.T, a *api.API) http.Handler
type testCase func(t *testing.T, bootstrap serverBootstrap)

func NewTestSuite() map[string]testCase {
	return map[string]testCase{
		"/api/fast": func(t *testing.T, sb serverBootstrap) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			v := domain.NewValuer()
			s := domain.NewMockSleeper(ctr)
			a := api.New(v, s)
			handler := sb(t, a)

			fx := httpTest.NewFixture(t, handler)
			defer fx.Finish()

			resp := fx.Request("GET", "/api/fast", nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			var body responses.Answer
			fx.Parse(resp.Body, &body)
			require.True(t, 0 < body.Value && body.Value < 101, fmt.Sprintf("value = %d", body.Value))
		},
		"/api/slow": func(t *testing.T, sb serverBootstrap) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			v := domain.NewValuer()
			s := domain.NewMockSleeper(ctr)
			s.EXPECT().Sleep(1000 * time.Millisecond)
			a := api.New(v, s)
			handler := sb(t, a)

			fx := httpTest.NewFixture(t, handler)
			defer fx.Finish()

			resp := fx.Request("GET", "/api/slow", nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			var body responses.Answer
			fx.Parse(resp.Body, &body)
			require.True(t, 0 < body.Value && body.Value < 101, fmt.Sprintf("value = %d", body.Value))
		},
		"/api/rand": func(t *testing.T, sb serverBootstrap) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			v := domain.NewValuer()
			s := domain.NewMockSleeper(ctr)
			s.EXPECT().SleepInterval(100*time.Millisecond, 1000*time.Millisecond)
			a := api.New(v, s)
			handler := sb(t, a)

			fx := httpTest.NewFixture(t, handler)
			defer fx.Finish()

			resp := fx.Request("GET", "/api/random", nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			var body responses.Answer
			fx.Parse(resp.Body, &body)
			require.True(t, 0 < body.Value && body.Value < 101, fmt.Sprintf("value = %d", body.Value))
		},
		"/api/rand (error case)": func(t *testing.T, sb serverBootstrap) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			v := domain.NewValuer()
			s := domain.NewMockSleeper(ctr)
			s.EXPECT().SleepInterval(100*time.Millisecond, 1000*time.Millisecond).Return(errors.New(""))
			a := api.New(v, s)
			handler := sb(t, a)

			fx := httpTest.NewFixture(t, handler)
			defer fx.Finish()

			resp := fx.Request("GET", "/api/random", nil)
			require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		},
		"/api/never_ending": func(t *testing.T, sb serverBootstrap) {
			ctr := gomock.NewController(t)
			defer ctr.Finish()

			v := domain.NewValuer()
			s := domain.NewMockSleeper(ctr)
			s.EXPECT().LoopForever()
			a := api.New(v, s)
			handler := sb(t, a)

			fx := httpTest.NewFixture(t, handler)
			defer fx.Finish()

			resp := fx.Request("GET", "/api/never_ending", nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		},
	}
}
