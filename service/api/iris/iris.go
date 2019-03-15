package iris

import (
	"highload/http_adapters"
	"highload/service/api"

	"github.com/kataras/iris"
)

func New(a *api.API) *iris.Application {
	server := iris.Default()

	// TODO
	// Add necessary middleware (for logging info about each user, for example)

	server.Get("api/endpoint", httpAdapters.ToIris(a.Endpoint))

	return server
}
