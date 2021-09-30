package main

import (
	"SIMGLEPROXY/myhttp"
	"fmt"
	"net"
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

	request, err := myhttp.ParseHttpRequest(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", request)

}
