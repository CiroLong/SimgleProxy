package myhttp

func newResponse() Response {
	Res := new(Response)
	Res.Headers = make(headers)
	Res.Body = make([]byte, bufMaxSize)
	return *Res
}
