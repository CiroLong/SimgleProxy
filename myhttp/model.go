package myhttp

type headers map[string][]string
type Request struct {
	Method string // "GET"
	Url    string // 自己解析,还是不用*url.URL了
	Proto  string // "HTTP/1.0"

	Headers headers

	Body []byte
}

type Url struct {
	Router string
	Query  map[string][]string // key=value1, key = value2
}

type Responce struct {
	Proto      string
	StatusCode int
	Phrase     string

	Headers headers

	Body []byte
}

const bufMaxSize = 1024 * 8

func newRequest() Request {
	r := new(Request)
	r.Headers = make(headers)
	r.Body = make([]byte, bufMaxSize)
	return *r
}

func newResponce() Responce {
	Res := new(Responce)
	Res.Headers = make(headers)
	Res.Body = make([]byte, bufMaxSize)
	return *Res
}
