package proxy

import (
	"SIMGLEPROXY/myhttp"
	"net"
)

type Server interface {
	Serve(net.Conn, *myhttp.Request) error
}

type TargetServer struct {
	ProxySetHeader map[string][]string //Host $test.com
	ProxyPass      string              //http://127.0.0.1:7000
	LocationRouter string              //   `/api/v1`
}

type StaticServer struct {
	RemotePath      string //匹配的远端路径   "/static"
	LocalRoot       string //本地root地址    "/mnt/var/www/localhost"
	DefaultFilePath string //默认文件        "/mnt/var/www/localhost/index.html"
}

var TargetRegistered map[string]*TargetServer //利用LocationRouter作为键
var StaticRegistered map[string]*StaticServer //利用RemotePath做为键
