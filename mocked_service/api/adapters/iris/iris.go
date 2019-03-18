package iris

import (
	"highload/http_adapters"
	"highload/mocked_service/api"

	"github.com/kataras/iris"
)

func New(a *api.API) *iris.Application {
	server := iris.Default()

	server.Get("api/fast", httpAdapters.ToIris(a.Fast))
	server.Get("api/slow", httpAdapters.ToIris(a.Slow))
	server.Get("api/random", httpAdapters.ToIris(a.Random))
	server.Get("api/never_ending", httpAdapters.ToIris(a.NeverEnding))

	return server
}
