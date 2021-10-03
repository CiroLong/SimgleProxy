package myhttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"
)

func NewResponse() Response {
	Res := new(Response)
	Res.Proto = "http/1.1"
	Res.Headers = make(headers)
	Res.Body = make([]byte, bufMaxSize)
	Res.Headers["Server"] = append(Res.Headers["Server"], "proxyserver/ciro")
	Res.Headers["Date"] = append(Res.Headers["Date"], time.Now().String())
	return *Res
}

func (r *Response) NotFound() {
	r.StatusCode = 404
	r.Phrase = "Not Found"
	// Connection: keep-alive
	// Content-Length: 555
	// Content-Type: text/html
	// Date: Sun, 03 Oct 2021 12:00:41 GMT
	// Server: nginx/1.21.3

	r.Headers["Content-Type"] = append(r.Headers["Content-Type"], "text/html")
	buf := new(bytes.Buffer)
	buf.WriteString("<html>\r\n")
	buf.WriteString("<head><title>404 Not Found</title></head>\r\n")
	buf.WriteString("<body>\r\n")
	buf.WriteString("<center><h1>404 Not Found</h1></center>\r\n")
	buf.WriteString("<hr><center>nginx/1.21.3</center>\r\n")
	buf.WriteString("</body>\r\n")
	buf.WriteString("</html>")
	body, _ := ioutil.ReadAll(buf)
	r.Body = append(r.Body, body...)
	r.Headers["Content-Length"] = append(r.Headers["Content-Length"], fmt.Sprint(len(r.Body)))

}
