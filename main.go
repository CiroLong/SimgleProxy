package main

import (
	"SIMGLEPROXY/myhttp"
	"SIMGLEPROXY/proxy"
	"bufio"
	"fmt"
	"net"
)

func init() {
	proxy.Init()
}

func main() {
	start()

	proxy.Realse()
}

func start() {
	listenr, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	defer listenr.Close()

	for {
		conn, err := listenr.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		go Loggertest(conn)
	}
}

func Loggertest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	request, err := myhttp.ParseHttpRequest(reader)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}
	// fmt.Println("parse success")
	//fmt.Printf("%#v\n", request)

	s, ps, err := proxy.FindServer(&request)

	// _, werr := ps.AccessLogger.Fp.WriteString("test string")
	// if werr != nil {
	// 	fmt.Println("write err", werr)
	// } else {
	// 	fmt.Println("write success", werr)
	// }

	if err == nil {
		if err := s.Serve(conn, &request, &ps); err != nil {
			fmt.Println("Serve err", err)
			ps.ErrorLogger.PrintError(err)
		}
	} else {
		fmt.Println("find server error:", err)
		ps.ErrorLogger.PrintError(err)
	}
}

// func proxyTest(conn net.Conn) {
// 	defer conn.Close()

// 	reader := bufio.NewReader(conn)
// 	request, err := myhttp.ParseHttpRequest(reader)
// 	if err != nil {
// 		fmt.Println("Parse error:", err)
// 		return
// 	}
// 	fmt.Println("parse success")
// 	// fmt.Printf("%#v\n", request.Headers)

// 	s, _ := proxy.FindServer(request.UrlParsed.Router) //传router好点

// 	if s != nil {
// 		if err := s.Serve(conn, &request); err != nil {
// 			fmt.Println("Serve err", err)
// 		}
// 	} else {
// 		fmt.Println("s is nil !")
// 	}
// }

// func staticTest(conn net.Conn) {
// 	defer conn.Close()

// 	reader := bufio.NewReader(conn)
// 	request, err := myhttp.ParseHttpRequest(reader)
// 	if err != nil {
// 		fmt.Println("Parse error:", err)
// 		return
// 	}
// 	fmt.Println("parse success")
// 	// fmt.Printf("%#v\n", request.Headers)

// 	s, _ := proxy.FindServer(request.UrlParsed.Router) //传router好点

// 	if s != nil {
// 		if err := s.Serve(conn, &request); err != nil {
// 			fmt.Println("Serve err", err)
// 		}
// 	} else {
// 		fmt.Println("s is nil !")
// 	}
// }

// func handler(conn net.Conn) {
// 	defer conn.Close()

// 	// 暂定为直接转发给 127.0.0.1:8081
// 	// Host 为test.com

// 	reader := bufio.NewReader(conn)
// 	request, err := myhttp.ParseHttpRequest(reader)
// 	if err != nil {
// 		fmt.Println("Parse error:", err)
// 		return
// 	}
// 	fmt.Println("parse success")
// 	// fmt.Printf("%#v\n", request.Headers)

// 	//编辑请求
// 	newRequest := request.Copy()
// 	//newRequest.ChangeHost("127.0.0.1:8081")
// 	//？

// 	fmt.Println("copy success")
// 	// fmt.Printf("%#v\n", newRequest.Headers)

// 	proxy, err := net.Dial("tcp", "127.0.0.1:8081")
// 	if err != nil {
// 		fmt.Println("proxy connect error:", err)
// 		//返回 错误代码
// 		return
// 	}
// 	defer proxy.Close()

// 	fmt.Println("dial success")

// 	// _, err = io.Copy(conn, proxy) //好像这样就可以了
// 	// if err != nil {
// 	// 	fmt.Println("io copy error:", err)
// 	// 	return
// 	// }

// 	go func() {
// 		//应该在另一个goroutine中读取
// 		_, err := io.Copy(conn, proxy) //这个函数一直阻塞？
// 		if err != nil {
// 			fmt.Println("io copy err", err)
// 		}
// 	}()
// 	//myhttp.SendHttpRequest(os.Stdout, newRequest)
// 	//var responseData []byte = make([]byte, 1024*8)
// 	bufffff := new(bytes.Buffer)

// 	err = myhttp.SendHttpRequest(bufffff, newRequest)
// 	if err != nil {
// 		fmt.Println("send http Request error:", err)
// 		return
// 	}
// 	box, _ := ioutil.ReadAll(bufffff)
// 	proxy.Write(box)

// 	// fmt.Println("send success")

// 	//这一段要重构
// 	// responseReader := bufio.NewReader(proxy)

// 	// _, err = proxy.Read(responseData)
// 	// if err != nil && err != io.EOF {
// 	// 	fmt.Println("get response error", err)
// 	// 	fmt.Println(string(responseData))
// 	// 	return
// 	// }

// 	// _, err = conn.Write(responseData)
// 	// if err != nil {
// 	// 	fmt.Println("wirte back error", err)
// 	// 	return
// 	// }

// 	fmt.Println("done")

// 	select {}

// }
