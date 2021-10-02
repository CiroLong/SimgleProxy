package myhttp

type headers map[string][]string
type Request struct {
	Method    string // "GET"
	Url       string // 自己解析,还是不用*url.URL了
	UrlParsed Url
	Proto     string // "HTTP/1.0"

	Headers headers

	Body []byte
}

type Url struct {
	Router string
	Query  map[string][]string // key=value1, key = value2
}

type Response struct {
	Proto      string
	StatusCode int
	Phrase     string

	Headers headers

	Body []byte
}

const bufMaxSize = 1024 * 8
