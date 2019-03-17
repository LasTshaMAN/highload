package iris

import (
	"highload/http_adapters"
	"highload/service/api"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kataras/iris"
)

func New(a *api.API, middleware ...iris.Handler) *iris.Application {
	i := iris.Default()

	regMiddleware(i, middleware...)
	regAPI(i, a)
	regMetrics(i)

	return i
}

func regMiddleware(i *iris.Application, middleware ...iris.Handler) {
	i.Use(middleware...)
}

func regAPI(i *iris.Application, a *api.API) {
	i.Get("api/endpoint", httpAdapters.ToIris(a.Endpoint))
}

func regMetrics(i *iris.Application) {
	i.Get("/metrics", iris.FromStd(promhttp.Handler()))
}
