package httpAdapters

import (
	"fmt"
	"time"

	"github.com/kataras/iris/context"

	"github.com/kataras/iris"
)

func RunIris(i *iris.Application, port int) (*iris.Application, error) {
	e := make(chan error)
	go func() {
		e <- i.Run(iris.Addr(fmt.Sprintf(":%d", port)))
	}()
	select {
	case err := <-e:
		return nil, err
	case <-time.After(100 * time.Millisecond):
		return i, nil
	}
}

func ToIris(h handler) iris.Handler {
	return func(irisCtx context.Context) {
		ctx := newIrisCtx(irisCtx)
		h(ctx)
	}
}
