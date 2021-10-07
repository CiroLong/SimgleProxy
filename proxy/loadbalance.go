package proxy

import "strings"

func LoadProxyPass(proxyPass string) string {

	proxyPass_cut := strings.ReplaceAll(proxyPass, "http://", "")

	if servers, ok := UpStreamsRegistered[proxyPass_cut]; ok {
		rlt := servers.ServerName[servers.now]
		servers.now++
		if servers.now >= len(servers.ServerName) {
			servers.now = 0
		}
		rlt = strings.ReplaceAll(rlt, "http://", "")
		return rlt
	}

	return proxyPass
}
