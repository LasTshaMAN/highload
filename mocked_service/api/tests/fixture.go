package tests

import (
	"highload/http_test"
	"highload/mocked_service/api"
	"highload/mocked_service/domain"
	"testing"

	"github.com/golang/mock/gomock"
)

type fixture struct {
	*httpTest.Fixture
	ctr     *gomock.Controller
	sleeper *domain.MockSleeper
}

func newFixture(t *testing.T, sb serverBootstrap) *fixture {
	ctr := gomock.NewController(t)

	v := domain.NewValuer()
	s := domain.NewMockSleeper(ctr)
	a := api.New(v, s)
	handler := sb(t, a)

	return &fixture{
		Fixture: httpTest.NewFixture(t, handler),
		ctr:     ctr,
		sleeper: s,
	}
}
