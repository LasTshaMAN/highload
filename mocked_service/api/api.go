package api

import (
	"highload/mocked_service/api/responses"
	"highload/mocked_service/domain/domain"
	"time"
)

type API struct {
	v domain.Valuer
	s domain.Sleeper
}

func New(v domain.Valuer, s domain.Sleeper) *API {
	return &API{
		v: v,
		s: s,
	}
}

func (a *API) Fast(ctx Context) {
	response := responses.Answer{
		Value: a.v.Value(),
	}
	ctx.JSON(response)
}

func (a *API) Slow(ctx Context) {
	a.s.Sleep(1000 * time.Millisecond)
	response := responses.Answer{
		Value: a.v.Value(),
	}
	ctx.JSON(response)
}

func (a *API) Random(ctx Context) {
	if err := a.s.SleepInterval(100*time.Millisecond, 1000*time.Millisecond); err != nil {
		ctx.ServerError(err)
		return
	}
	response := responses.Answer{
		Value: a.v.Value(),
	}
	ctx.JSON(response)
}

func (a *API) NeverEnding(ctx Context) {
	a.s.LoopForever()
}
