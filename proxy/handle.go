package proxy

import (
	"SIMGLEPROXY/myhttp"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Init() {
	// TargetRegistered = make(map[string]*TargetServer)
	// StaticRegistered = make(map[string]*StaticServer)
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
		p.Locations = make(map[string]Server)
		p.ServerName = server.ServerName
		p.AccessLogPath = server.AccessLogPath
		p.ErrorLogPath = server.ErrorLogPath
		p.ListenPort = strconv.Itoa(server.Listen)
		for _, location := range server.Locations {
			if location.IsStatic {
				ss := NewStaticServer(StaticServer{
					RemotePath:      location.Router,
					LocalRoot:       location.Root,
					DefaultFilePath: location.Root + "/" + location.Index,
				})
				//p.Locations = append(p.Locations, ss)
				p.Locations[ss.RemotePath] = ss
			} else {
				ts := NewTargetServer(TargetServer{
					ProxyPass:      strings.ReplaceAll(location.ProxyPass, "http://", ""),
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
				//p.Locations = append(p.Locations, ts)
				p.Locations[ts.LocationRouter] = ts
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
func FindServer(req *myhttp.Request) (Server, error) {

	//首先应该根据 Host 头部确定对应 server_name?
	//应该叫server_name
	hosts, ok := req.Headers["Host"]
	if !ok {
		return NewTargetServer(TargetServer{}), errors.New("no Host Header")
	}
	host := hosts[0] //简化

	ps, ok := ProxyServerRegistered[host]
	if !ok {
		return NewTargetServer(TargetServer{}), errors.New("no such proxy server")
	}

	//
	// 这里可能要对日志做处理
	//

	//然后根据 router 确定对应的location  (之后考虑通配符)
	//回头写个函数
	for r, s := range ps.Locations { //woc完蛋，匹配不了， proxyServer要改
		if req.UrlParsed.Router == r { //暂时用等于匹配,还需要考虑静态文件匹配,通配符匹配的问题
			//得到Server实例,返回
			server := s
			return server, nil
		}
	}

	return NewTargetServer(TargetServer{}), errors.New("no matching location")
}

//原来的没用了，先放这
// func FindServer(url string) (server Server, isStatic bool) {
// 	for key, val := range StaticRegistered {
// 		if strings.HasPrefix(url, key) {
// 			isStatic = true
// 			temp := val //还是写一下好
// 			server = temp
// 			break
// 		}
// 	}

// 	if isStatic {
// 		return
// 	}

// 	for key, val := range TargetRegistered {
// 		if strings.HasPrefix(url, key) {
// 			temp := val
// 			server = temp
// 			break
// 		}
// 	}
// 	return
// }
