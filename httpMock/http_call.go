package httpMock

type httpCall struct {
	URL string

	RespStatus int
	RespBody   []byte

	Times int
}

func (mc *httpCall) uniqueId() uniqueId {
	return uniqueId{
		URL: mc.URL,
	}
}

type uniqueId struct {
	URL string
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
		URL:        b.url,
		RespStatus: b.respStatus,
		RespBody:   b.respBody,
		Times:      b.times,
	}
}
