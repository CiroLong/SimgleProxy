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

func (s *TargetServer) Serve(conn net.Conn, request *myhttp.Request) error {
	//fmt.Printf("%#v\n", *s)

	//编辑请求
	newRequest := request.Copy()
	newRequest.ChangeHost(s.ProxyPass)
	//？

	fmt.Println("copy success")
	// fmt.Printf("%#v\n", newRequest.Headers)

	proxy, err := net.Dial("tcp", s.ProxyPass)
	if err != nil {
		fmt.Println("proxy connect error:", err)
		//返回 错误代码
		return err
	}
	defer proxy.Close()

	fmt.Println("dial success")

	dataBuf := make([]byte, 1024*8)
	ok := make(chan struct{})
	go func() {
		proxy.Read(dataBuf)
		ok <- struct{}{}
	}()
	// go func() {
	// 	//应该在另一个goroutine中读取

	// 	_, err := io.Copy(conn, proxy) //这个函数一直阻塞？
	// 	if err != nil {
	// 		fmt.Println("io copy err", err)
	// 	}
	// }()

	bufffff := new(bytes.Buffer)

	err = myhttp.SendHttpRequest(bufffff, newRequest)
	if err != nil {
		fmt.Println("send http Request error:", err)
		return err
	}
	box, _ := ioutil.ReadAll(bufffff)
	proxy.Write(box)

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

	<-ok
	//fmt.Println(string(dataBuf)) //为什么这个打印了下面的Write写入不成功，而且err还是nil, 好吧是我的postman好像出问题了
	n, err := conn.Write(dataBuf)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("write 0 bytes?")
	}
	return nil
}
