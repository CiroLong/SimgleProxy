package myhttp

type headers map[string][]string
type Request struct {
	Method string // "GET"
	Url    string // 自己解析,还是不用*url.URL了
	Proto  string // "HTTP/1.0"

	Headers headers
}

const bufMaxSize = 1024 * 8

func newRequest() Request {
	r := new(Request)
	r.Headers = make(headers)
	return *r
}
