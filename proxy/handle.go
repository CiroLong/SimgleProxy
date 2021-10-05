package proxy

import (
	"SIMGLEPROXY/myhttp"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Init() {
	TargetRegistered = make(map[string]*TargetServer)
	StaticRegistered = make(map[string]*StaticServer)
	ProxyServerRegistered = make(map[string]ProxyServer)

	//这里应该用config读取配置文件
	c, err := LoadConfig()
	if err != nil {
		fmt.Println("LoadConfig err:", err)
		return
	}
	fmt.Printf("%#v\n\n", c)

	for _, server := range c.Servers {
		p := new(ProxyServer)
		p.ServerName = server.ServerName
		p.AccessLogPath = server.AccessLogPath
		p.ErrorLogPath = server.ErrorLogPath
		p.ListenPort = strconv.Itoa(server.Listen)
		for _, location := range server.Locations {
			if location.IsStatic {
				ss := NewStaticServer(StaticServer{
					RemotePath:      location.Router,
					LocalRoot:       location.Root,
					DefaultFilePath: location.Index,
				})
				p.Locations = append(p.Locations, ss)
			} else {
				ts := NewTargetServer(TargetServer{
					ProxyPass:      location.ProxyPass,
					LocationRouter: location.Router,
				})
				//难过
				ts.ProxySetHeader = make(map[string][]string)
				x := strings.Split(location.ProxySetHeader, "&")
				for _, xx := range x {
					y := strings.Split(xx, "=")
					if len(y) != 2 {
						fmt.Println("config ProxySetHeader error")
						return
					}
					ts.ProxySetHeader[y[0]] = append(ts.ProxySetHeader[y[0]], y[1])
				}
				p.Locations = append(p.Locations, ts)
			}
		}

		ProxyServerRegistered[p.ServerName] = *p
	}

	// fmt.Printf("%#v\n\n", ProxyServerRegistered)

	//test
	// for _, s := range ProxyServerRegistered {
	// 	for _, server := range s.Locations {
	// 		fmt.Printf("%#v\n", server)
	// 	}
	// }

	//test ok, right

	// TargetRegistered["/api/v1"] = NewTargetServer(TargetServer{
	// 	ProxySetHeader: (map[string][]string{
	// 		"Host": {"test.com"},
	// 	}),
	// 	ProxyPass:      "127.0.0.1:8081",
	// 	LocationRouter: "/api/v1",
	// })

	// StaticRegistered["/static"] = NewStaticServer(StaticServer{
	// 	RemotePath:      "/static",
	// 	LocalRoot:       "/home/ciro/mydocument/localhost",
	// 	DefaultFilePath: "/home/ciro/mydocument/localhost/index.html",
	// })
}

//由于更改了全局对象存放方式,FindServer函数要大改
//要不传入request指针, 在函数内部确定对象好了
func FindServer2(req *myhttp.Request) (server Server, isStatic bool, err error) {

	//首先应该根据 Host 头部确定对应 server_name?
	//应该叫server_name
	hosts, ok := req.Headers["Host"]
	if !ok {
		return NewTargetServer(""), false, errors.New("No Host Header")
	}
	host := hosts[0] //简化

	ps, ok := ProxyServerRegistered[host]
	if !ok {
		return NewTargetServer(""), false, errors.New("No such proxy server")
	}

	//
	// 这里可能要对日志做处理
	//

	//然后根据 router 确定对应的location  (之后考虑通配符)
	//回头写个函数
	for _, server := range ps.Locations { //woc完蛋，匹配不了， proxyServer要改
		//if req.UrlParsed.Router ==
	}

	//得到Server实例,返回

	return NewTargetServer(""), false, nil
}

//原来的没用了，先放这
func FindServer(url string) (server Server, isStatic bool) {
	for key, val := range StaticRegistered {
		if strings.HasPrefix(url, key) {
			isStatic = true
			temp := val //还是写一下好
			server = temp
			break
		}
	}

	if isStatic {
		return
	}

	for key, val := range TargetRegistered {
		if strings.HasPrefix(url, key) {
			temp := val
			server = temp
			break
		}
	}
	return
}
