package api

type Handler func(Context)

type Context interface {
	JSON(response interface{})
	ServerError(err error)
}
