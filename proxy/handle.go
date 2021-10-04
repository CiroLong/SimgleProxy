package proxy

import (
	"fmt"
	"strings"
)

func Init() {
	TargetRegistered = make(map[string]*TargetServer)
	StaticRegistered = make(map[string]*StaticServer)

	//这里应该用config读取配置文件
	c, err := LoadConfig()
	if err != nil {
		fmt.Println("LoadConfig err:", err)
		return
	}
	fmt.Printf("%#v\n", c)

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
	//?
}

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
