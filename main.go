package main

import (
	"SIMGLEPROXY/myhttp"
	"SIMGLEPROXY/proxy"
	"bufio"
	"fmt"
	"net"
	"sync"
)

func init() {
	proxy.Init()
}

func main() {
	defer proxy.Release()
	start()
}

func start() {
	portMap := make(map[string]bool)
	portMapLock := new(sync.Mutex)

	for _, ps := range proxy.ProxyServerRegistered {
		nps := ps
		go func() {

			portMapLock.Lock()
			if _, ok := portMap[nps.ListenPort]; ok {
				portMapLock.Unlock()
				return
			}
			portMap[nps.ListenPort] = true
			port := ":" + nps.ListenPort
			listenr, err := net.Listen("tcp", port)
			if err != nil {
				fmt.Println("listen error:", err)
				return
			}
			defer listenr.Close()

			portMapLock.Unlock()

			for {
				conn, err := listenr.Accept()
				if err != nil {
					fmt.Println("accept error:", err)
					break
				}
				go test(conn)
			}
		}()
	}
	select {}
}

func test(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	request, err := myhttp.ParseHttpRequest(reader)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	s, ps, err := proxy.FindServer(&request)

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
