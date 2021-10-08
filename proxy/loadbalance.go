package proxy

import "strings"

func LoadProxyPass(proxyPass string) string {

	proxyPass_cut := strings.ReplaceAll(proxyPass, "http://", "")

	servers, ok := UpStreamsRegistered[proxyPass_cut]
	if ok {
		if servers.Mode == "PR" { //轮询
			rlt := servers.ServerName[servers.now]
			servers.now++
			if servers.now >= len(servers.ServerName) {
				servers.now = 0
			}
			rlt = strings.ReplaceAll(rlt, "http://", "")
			return rlt
		} else if servers.Mode == "WRR" { //加权轮询
			all := 0
			get := false
			var rlt string
			for i, w := range servers.Weight {
				if w <= 0 {
					break
				}
				all += w
				if servers.now < all && !get {
					rlt = servers.ServerName[i]
					servers.now++
					get = true
					rlt = strings.ReplaceAll(rlt, "http://", "")
				}
			}
			if servers.now >= all {
				servers.now = 0
			}
			return rlt
		}
	}

	return proxyPass
}
