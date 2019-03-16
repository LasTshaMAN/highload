package httpMock

type httpCall struct {
	url string

	respStatus int
	respBody   []byte

	times int
}

func (mc *httpCall) uniqueId() uniqueId {
	return uniqueId{
		url: mc.url,
	}
}

type uniqueId struct {
	url string
}

type httpCallBuilder struct {
	c *FrontEnd

	url string

	respStatus int
	respBody   []byte

	times int
}

func (b *httpCallBuilder) Return(status int, body []byte) *httpCallBuilder {
	b.respStatus = status
	b.respBody = body

	b.c.register(b.build())

	return b
}

func (b *httpCallBuilder) Times(t int) {
	b.times = t

	b.c.reRegister(b.build())
}

func (b *httpCallBuilder) build() httpCall {
	return httpCall{
		url:        b.url,
		respStatus: b.respStatus,
		respBody:   b.respBody,
		times:      b.times,
	}
}
