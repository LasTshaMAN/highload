package api

import (
	"highload/http_adapters"
	"highload/service/api/responses"
	"highload/service/domain"
)

type API struct {
	avg domain.Avg
}

func New(avg domain.Avg) *API {
	return &API{
		avg: avg,
	}
}

func (a *API) Endpoint(ctx httpAdapters.Context) {
	value, err := a.avg.Value()
	if err != nil {
		ctx.ServerError(err)
		return
	}
	response := responses.Answer{
		Value: value,
	}
	ctx.JSON(response)
}
