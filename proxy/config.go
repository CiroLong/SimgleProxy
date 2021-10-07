package proxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

//这是用来读取配置文件的包
//包括模型和函数
//以及对应的TargetServer和StaticServer构建

type Location struct {
	Router         string `json:"router"`
	ProxySetHeader string `json:"proxy_set_header"`
	ProxyPass      string `json:"proxy_pass"`

	Root  string `json:"root"`
	Index string `json:"index"`

	IsStatic bool `json:"isstatic"`
}

type Server_config struct {
	Listen        int        `json:"listen"` //listen port
	ServerName    string     `json:"server_name"`
	Locations     []Location `json:"locations"`
	ErrorLogPath  string     `json:"error_log"`
	AccessLogPath string     `json:"access_log"`
}

type UpStreamCofig struct {
	ProxyPass  string   `json:"proxy_pass"`
	ServerName []string `json:"server"`
}

type Config struct {
	UpStreams []UpStreamCofig `json:"upstreams"`
	Servers   []Server_config `json:"servers"`
}

func LoadConfig() (Config, error) {
	c := new(Config)
	fp, err := os.Open("./config.json")
	if err != nil {
		return *c, err
	}
	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return *c, err
	}

	if err := json.Unmarshal(data, &c); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}

	return *c, nil
}
