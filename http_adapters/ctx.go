package httpAdapters

type handler func(Context)

type Context interface {
	JSON(response interface{})
	ServerError(err error)
}
