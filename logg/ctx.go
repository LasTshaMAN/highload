package logg

import (
	"context"

	"github.com/kataras/iris"
)

const ctxKey string = "log"

// TODO
// Test it !
func FromContext(ctx iris.Context) Logger {
	res, ok := ctx.Values().Get(ctxKey).(Logger)
	if !ok {
		return fallback{}
	}
	return res
}

// TODO
// Test it !
func ToContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxKey, l)
}
