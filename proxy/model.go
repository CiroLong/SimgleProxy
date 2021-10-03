package proxy

type TargetServer struct {
	Proxy_set_header map[string][]string //Host $test.com
	Proxy_pass       string              //http://127.0.0.1:7000
}

type StaticServer struct {
	RemotePath      string //匹配的远端路径   "/static"
	LocalRoot       string //本地root地址    "/mnt/var/www/localhost"
	DefaultFileName string //默认文件        "/mnt/var/www/localhost/index.html"
}
