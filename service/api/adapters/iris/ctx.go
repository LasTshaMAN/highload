package iris

import "github.com/kataras/iris"

type ContextImpl struct {
	irisCtx iris.Context
}

func newContext(irisCtx iris.Context) *ContextImpl {
	return &ContextImpl{
		irisCtx: irisCtx,
	}
}

// TODO
// Test it!
func (c *ContextImpl) JSON(response interface{}) {
	if _, err := c.irisCtx.JSON(response); err != nil {
		// TODO
		// Handle err
		//logg.FromContext(ctx).Errorf(err.Error())
	}
}
