package proxy

import (
	"SIMGLEPROXY/myhttp"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func NewStaticServer(info interface{}) *StaticServer {
	s := new(StaticServer)
	*s = info.(StaticServer)
	return s
}

//传入处理后的字符串,如 `/img/test.png`
func (s *StaticServer) GetFile(path string) (*os.File, error) {
	var fullPath string
	if path == "" {
		fullPath = s.DefaultFilePath
	} else {
		fullPath = s.LocalRoot + path
	}
	fmt.Println("fullpath:", fullPath)
	fp, err := os.Open(fullPath)
	return fp, err
}

func (s *StaticServer) Serve(conn net.Conn, request *myhttp.Request) error {

	path := request.UrlParsed.Router[:]
	path = strings.Replace(path, s.RemotePath, "", 1)

	fp, err := s.GetFile(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	resp := myhttp.NewResponse()
	body, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}
	resp.Body = body[:]
	resp.Headers["Content-Length"] = append(resp.Headers["Content-Length"], fmt.Sprint(len(resp.Body)))
	resp.Headers["Content-Type"] = append(resp.Headers["Content-Type"], "text/html") //得找个办法判断一下文件类型
	resp.StatusCode = 200
	resp.Phrase = "OK"

	err = myhttp.SendHttpResponse(conn, resp)
	if err != nil {
		return err
	}
	fmt.Println("done")
	return nil
}
