package domain_test

import (
	"crypto/tls"
	"encoding/json"
	"highload/httpMock"
	"highload/service/domain"
	"net/http"
	"testing"

	"golang.org/x/net/http2"

	"github.com/stretchr/testify/require"
)

func TestSequentialAvg_Value(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		vFast := domain.Answer{Value: 1}
		nFast := 5
		fx.Mock("/api/fast").Return(http.StatusOK, fx.toJSON(vFast)).Times(nFast)
		vSlow := domain.Answer{Value: 100}
		nSlow := 5
		fx.Mock("/api/slow").Return(http.StatusOK, fx.toJSON(vSlow)).Times(nSlow)
		vRandom := domain.Answer{Value: 50}
		nRandom := 20
		fx.Mock("/api/random").Return(http.StatusOK, fx.toJSON(vRandom)).Times(nRandom)

		client := &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		avg := domain.NewSequentialAvg(fx.Host(), client)
		act, err := avg.Value()
		require.NoError(t, err)
		exp := (vFast.Value*nFast + vSlow.Value*nSlow + vRandom.Value*nRandom) / (nFast + nSlow + nRandom)
		require.Equal(t, exp, act)
	})
}

func TestConcurrentAvg_Value(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		vFast := domain.Answer{Value: 1}
		nFast := 5
		fx.Mock("/api/fast").Return(http.StatusOK, fx.toJSON(vFast)).Times(nFast)
		vSlow := domain.Answer{Value: 100}
		nSlow := 5
		fx.Mock("/api/slow").Return(http.StatusOK, fx.toJSON(vSlow)).Times(nSlow)
		vRandom := domain.Answer{Value: 50}
		nRandom := 20
		fx.Mock("/api/random").Return(http.StatusOK, fx.toJSON(vRandom)).Times(nRandom)

		client := &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		avg := domain.NewConcurrentAvg(fx.Host(), client)
		act, err := avg.Value()
		require.NoError(t, err)
		exp := (vFast.Value*nFast + vSlow.Value*nSlow + vRandom.Value*nRandom) / (nFast + nSlow + nRandom)
		require.Equal(t, exp, act)
	})
}

type fixture struct {
	*httpMock.FrontEnd

	t *testing.T
}

func newFixture(t *testing.T) *fixture {
	return &fixture{
		FrontEnd: httpMock.New(t),
		t:        t,
	}
}

func (fx *fixture) toJSON(v interface{}) []byte {
	result, err := json.Marshal(v)
	require.NoError(fx.t, err)
	return result
}
