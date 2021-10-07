package proxy

import (
	"SIMGLEPROXY/logger"
	"SIMGLEPROXY/myhttp"
	"net"
)

type Server interface {
	Serve(net.Conn, *myhttp.Request, *ProxyServer) error
}

type TargetServer struct {
	LocationRouter string              //   `/api/v1`
	ProxySetHeader map[string][]string //Host $test.com
	ProxyPass      string              //    http://127.0.0.1:7000

}

type StaticServer struct {
	RemotePath      string //匹配的远端路径   "/static"
	LocalRoot       string //本地root地址    "/mnt/var/www/localhost"
	DefaultFilePath string //默认文件        "/mnt/var/www/localhost/index.html"
}

// var TargetRegistered map[string]*TargetServer //利用LocationRouter作为键
// //好像该用server_name - > 29

// var StaticRegistered map[string]*StaticServer //利用RemotePath做为键

//换个方案
type ProxyServer struct {
	ListenPort    string
	ServerName    string
	Locations     map[string]Server // [=|^~|~|~*|@] path    (利用RemotePath或者LocationRouter做为键(x))
	ErrorLogPath  string
	AccessLogPath string
	ErrorLogger   logger.Logger
	AccessLogger  logger.Logger
}

var ProxyServerRegistered map[string]*ProxyServer //用server_name标识？

//以下为负载均衡搭建的model

type UpStream struct {
	ProxyPass  string
	ServerName []string
	now        int
}

var UpStreamsRegistered map[string]*UpStream // 暂时就写个轮循吧
