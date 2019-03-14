package iris

import (
	"fmt"
	"highload/service/api"
	"time"

	"github.com/kataras/iris/context"

	"github.com/kataras/iris"
)

// TODO
// avoid name collisions with `echo`

func New(a *api.API) *iris.Application {
	server := iris.Default()

	// TODO
	// Add necessary middleware (for logging info about each user, for example)

	server.Get("api/fast", toIris(a.Fast))
	server.Get("api/slow", toIris(a.Slow))
	server.Get("api/random", toIris(a.Random))
	server.Get("api/never_ending", toIris(a.NeverEnding))

	return server
}

func Run(a *api.API, port int) (*iris.Application, error) {
	server := New(a)

	e := make(chan error)
	go func() {
		e <- server.Run(iris.Addr(fmt.Sprintf(":%d", port)))
	}()
	select {
	case err := <-e:
		return nil, err
	case <-time.After(100 * time.Millisecond):
		return server, nil
	}
}

func toIris(handler api.Handler) iris.Handler {
	return func(irisCtx context.Context) {
		ctx := newContext(irisCtx)
		handler(ctx)
	}
}
