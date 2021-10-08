package proxy

import (
	"SIMGLEPROXY/myhttp"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Init() {
	// TargetRegistered = make(map[string]*TargetServer)
	// StaticRegistered = make(map[string]*StaticServer)
	ProxyServerRegistered = make(map[string]*ProxyServer)

	UpStreamsRegistered = make(map[string]*UpStream)

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
		p.AccessLogger.Fp, err = os.OpenFile(p.AccessLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("open AccessLogger fail:", err)
			os.Exit(1)
		}
		p.ErrorLogger.Fp, err = os.OpenFile(p.ErrorLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("open ErrorLogger fail:", err)
			os.Exit(1)
		}

		ProxyServerRegistered[p.ServerName] = p
	}

	for _, upsrteamConfig := range c.UpStreams {
		u := new(UpStream)
		u.Weight = make([]int, 0)
		u.ServerName = make([]string, 0)

		u.ProxyPass = upsrteamConfig.ProxyPass
		u.ServerName = append(u.ServerName, upsrteamConfig.ServerName...)

		u.Weight = append(u.Weight, upsrteamConfig.Weight...)
		u.Mode = upsrteamConfig.Mode
		UpStreamsRegistered[u.ProxyPass] = u
	}
}

func Release() {
	for _, ps := range ProxyServerRegistered {
		ps.AccessLogger.Fp.Close()
		ps.ErrorLogger.Fp.Close()
	}
}

//由于更改了全局对象存放方式,FindServer函数要大改
//要不传入request指针, 在函数内部确定对象好了
func FindServer(req *myhttp.Request) (Server, ProxyServer, error) {

	//首先应该根据 Host 头部确定对应 server_name?
	//应该叫server_name
	hosts, ok := req.Headers["Host"]
	if !ok {
		return NewTargetServer(TargetServer{}), ProxyServer{}, errors.New("no Host Header")
	}
	host := hosts[0]                         //简化
	nhost := host[:strings.Index(host, ":")] // 去除端口号

	ps, ok := ProxyServerRegistered[nhost]
	if !ok {
		return NewTargetServer(TargetServer{}), *ps, errors.New("no such proxy server")
	}

	//然后根据 router 确定对应的location  (之后考虑通配符)
	//回头写个函数
	// for r, s := range ps.Locations { //woc完蛋，匹配不了， proxyServer要改

	// 	if req.UrlParsed.Router == r { //暂时用等于匹配,还需要考虑静态文件匹配,通配符匹配的问题
	// 		//得到Server实例,返回
	// 		server := s
	// 		return server, nil
	// 	}
	// }
	server, err := ps.LocationRouterMatch(req)
	if err != nil {
		return server, *ps, err
	}

	return server, *ps, nil
}

func (ps *ProxyServer) LocationRouterMatch(req *myhttp.Request) (Server, error) {
	//[=|^~|~|~*|@] path
	//var server Server

	url := req.UrlParsed
	var hasNormal bool

	for fpath, s := range ps.Locations {
		pathSlice := strings.Split(fpath, " ")
		if len(pathSlice) == 2 {
			if pathSlice[0] == "=" { // 精确匹配 =
				if url.Router == pathSlice[1] {
					server := s
					return server, nil
				}
			}
		}
	}

	var max int
	var pathPrefix string
	for fpath := range ps.Locations {
		pathSlice := strings.Split(fpath, " ")
		if len(pathSlice) == 2 {
			if pathSlice[0] == "^~" { // 前缀匹配 (非正则匹配 ^~ 返回匹配项多的？
				if strings.HasPrefix(url.Router, pathSlice[1]) {
					numRouter := strings.Split(pathSlice[1], "/")
					if len(numRouter) > max {
						max = len(numRouter)
						pathPrefix = fpath
					}
				}
			}
		}
	}
	if max > 0 {
		return ps.Locations[pathPrefix], nil
	}

	for fpath, s := range ps.Locations {
		pathSlice := strings.Split(fpath, " ")
		if len(pathSlice) == 2 {
			if pathSlice[0] == "~" { // 正则匹配 ~ 和 ~*
				//区分大小写
				match, _ := regexp.MatchString(pathSlice[1], url.Router)
				if match {
					server := s
					return server, nil
				}
			} else if pathSlice[0] == "~*" {
				//不区分大小写
				r := strings.ToLower(url.Router)
				p := strings.ToLower(pathSlice[1])
				match, _ := regexp.MatchString(p, r)
				if match {
					server := s
					return server, nil
				}
			}
		}
	}

	for fpath, s := range ps.Locations {
		pathSlice := strings.Split(fpath, " ")
		if len(pathSlice) == 1 { // 普通前缀匹配
			if pathSlice[0] == "/" {
				hasNormal = true
				continue
			}
			if strings.HasPrefix(url.Router, pathSlice[0]) {
				server := s
				return server, nil
			}
		}
	}

	if hasNormal {
		return ps.Locations["/"], nil
	} else {
		host := req.Headers["Host"]
		out := host[0] + req.Url + " (no matched location)"

		return nil, errors.New(out)
	}
}
