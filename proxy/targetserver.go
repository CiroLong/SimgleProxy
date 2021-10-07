package proxy

import (
	"SIMGLEPROXY/myhttp"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
)

func NewTargetServer(info interface{}) *TargetServer {
	s := new(TargetServer)
	*s = info.(TargetServer)
	return s
}

func (s *TargetServer) Serve(conn net.Conn, request *myhttp.Request, ps *ProxyServer) error {
	//编辑请求
	newRequest := request.Copy()
	newRequest.ChangeHost(s.ProxyPass)

	realProxyPass := LoadProxyPass(s.ProxyPass)

	proxy, err := net.Dial("tcp", realProxyPass)
	if err != nil {
		fmt.Println("proxy connect error:", err)
		//返回 错误代码
		return err
	}
	defer proxy.Close()

	dataBuf := make([]byte, 1024*8)
	ok := make(chan struct{})
	go func() {
		proxy.Read(dataBuf)
		ok <- struct{}{}
	}()

	bufffff := new(bytes.Buffer)

	err = myhttp.SendHttpRequest(bufffff, newRequest)
	if err != nil {
		fmt.Println("send http Request error:", err)
		return err
	}
	box, _ := ioutil.ReadAll(bufffff)
	proxy.Write(box)

	<-ok
	//fmt.Println(string(dataBuf)) //为什么这个打印了下面的Write写入不成功，而且err还是nil, 好吧是我的postman好像出问题了

	n, err := conn.Write(dataBuf)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("write 0 bytes?")
	}

	//没写response读取==
	dataBufff := new(bytes.Buffer)
	dataBufff.Write(dataBuf)
	response, _ := myhttp.ParseHttpResponse(dataBufff)

	ps.AccessLogger.PrintAccess(request, &response)
	return nil
}
