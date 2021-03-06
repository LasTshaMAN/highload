package tests

import (
	"highload/http_test"
	"highload/service/api"
	"highload/service/domain"
	"testing"

	"github.com/golang/mock/gomock"
)

type fixture struct {
	*httpTest.Fixture
	ctr *gomock.Controller
	avg *domain.MockAvg
}

func newFixture(t *testing.T, sb serverBootstrap) *fixture {
	ctr := gomock.NewController(t)

	avg := domain.NewMockAvg(ctr)
	a := api.New(avg)
	handler := sb(t, a)

	return &fixture{
		Fixture: httpTest.NewFixture(t, handler),
		ctr:     ctr,
		avg:     avg,
	}
}
