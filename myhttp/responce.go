package myhttp

func newResponce() Responce {
	Res := new(Responce)
	Res.Headers = make(headers)
	Res.Body = make([]byte, bufMaxSize)
	return *Res
}
