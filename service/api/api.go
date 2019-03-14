package api

import "highload/service/api/responses"

type API struct {
}

func New() *API {
	return &API{}
}

func (a *API) Fast(ctx Context) {
	response := responses.Answer{
		Value: 1,
	}
	ctx.JSON(response)
}

func (a *API) Slow(ctx Context) {
}

func (a *API) Random(ctx Context) {
}

func (a *API) NeverEnding(ctx Context) {
}
