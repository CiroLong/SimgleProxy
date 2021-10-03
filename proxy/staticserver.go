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
	// resp.NotFound()
	// myhttp.SendHttpResponse(os.Stdout, resp)
	// err = myhttp.SendHttpResponse(conn, resp)
	// if err != nil {
	// 	return err
	// }
	body, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}
	resp.Body = body[:]
	resp.Headers["Content-Length"] = append(resp.Headers["Content-Length"], fmt.Sprint(len(resp.Body)))
	resp.Headers["Content-Type"] = append(resp.Headers["Content-Type"], "text/html")
	resp.StatusCode = 200
	resp.Phrase = "OK"

	err = myhttp.SendHttpResponse(conn, resp)
	if err != nil {
		return err
	}

	//myhttp.SendHttpRequest(os.Stdout, newRequest)
	//var responseData []byte = make([]byte, 1024*8)
	// bufffff := new(bytes.Buffer)

	// err = myhttp.SendHttpRequest(bufffff, newRequest)
	// if err != nil {
	// 	fmt.Println("send http Request error:", err)
	// 	return
	// }
	// box, _ := ioutil.ReadAll(bufffff)
	// proxy.Write(box)

	// fmt.Println("send success")

	//这一段要重构
	// responseReader := bufio.NewReader(proxy)

	// _, err = proxy.Read(responseData)
	// if err != nil && err != io.EOF {
	// 	fmt.Println("get response error", err)
	// 	fmt.Println(string(responseData))
	// 	return
	// }

	// _, err = conn.Write(responseData)
	// if err != nil {
	// 	fmt.Println("wirte back error", err)
	// 	return
	// }
	fmt.Println("done")
	return nil
}
