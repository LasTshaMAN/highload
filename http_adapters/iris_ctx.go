package httpAdapters

import "github.com/kataras/iris"

type IrisCtx struct {
	irisCtx iris.Context
}

func newIrisCtx(irisCtx iris.Context) *IrisCtx {
	return &IrisCtx{
		irisCtx: irisCtx,
	}
}

// TODO
// Test it!
func (c *IrisCtx) JSON(response interface{}) {
	if _, err := c.irisCtx.JSON(response); err != nil {
		// TODO
		// Handle err
		//logg.FromContext(ctx).Errorf(err.Error())
	}
}

// TODO
// Test it!
func (c *IrisCtx) ServerError(err error) {
	c.irisCtx.Values().Set("error", err)
	c.irisCtx.StatusCode(iris.StatusInternalServerError)
}
