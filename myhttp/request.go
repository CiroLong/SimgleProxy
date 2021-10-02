package myhttp

func newRequest() Request {
	r := new(Request)
	r.Headers = make(headers)
	r.Body = make([]byte, bufMaxSize)
	return *r
}

func (r *Request) Copy() Request {
	copy := newRequest()
	copy.Method = r.Method
	copy.Url = r.Url
	copy.Proto = r.Proto
	copy.Body = r.Body[:]
	for k, v := range r.Headers {
		newv := v[:]
		copy.Headers[k] = append(copy.Headers[k], newv...)
	}

	return copy
}

func (r *Request) ChangeHost(newHost ...string) error {
	delete(r.Headers, "Host")
	r.Headers["Host"] = newHost[:] //?
	return nil
}
