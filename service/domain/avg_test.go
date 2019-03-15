package domain

import (
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAvgImpl_Value(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		fx := newFixture()
		defer fx.finish()

		vFast := 1
		nFast := 5
		fx.mock("/api/fast").value(vFast).times(nFast)
		vSlow := 100
		nSlow := 5
		fx.mock("/api/slow").value(vSlow).times(nSlow)
		vRandom := 50
		nRandom := 20
		fx.mock("/api/random").value(vRandom).times(nRandom)

		avg := NewAvg("", fx.client)
		act, err := avg.Value()

		exp := (vFast*nFast + vSlow*nSlow + vRandom*nRandom) / (nFast + nSlow + nRandom)
		require.NoError(t, err)
		require.Equal(t, exp, act)
	})
}

type fixture struct {
	client *http.Client
}

func newFixture() *fixture {
	client := &http.Client{}
	gock.InterceptClient(client)

	return &fixture{
		client: client,
	}
}

func (fx *fixture) finish() {
	gock.Off()
}

func (fx *fixture) mock(url string) *httpMock {
	return &httpMock{url: url}
}

type httpMock struct {
	url string
	v   int
}

func (hm *httpMock) value(v int) *httpMock {
	hm.v = v
	return hm
}

func (hm *httpMock) times(callCount int) {
	for i := 0; i < callCount; i++ {
		gock.New("").
			Get(hm.url).
			Reply(200).
			JSON(map[string]int{"value": hm.v})
	}
	return
}
