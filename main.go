package main

import (
	"SIMGLEPROXY/myhttp"
	"bufio"
	"fmt"
	"net"
	"os"
)

func init() {

}

func main() {

	start()
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
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()

	// 暂定为直接转发给 127.0.0.1:8081
	// Host 为test.com

	reader := bufio.NewReader(conn)
	request, err := myhttp.ParseHttpRequest(reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	myhttp.SendHttpRequest(os.Stdout, request)

	select {}

}
