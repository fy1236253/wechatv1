package http

import (
	"strings"
	"net/http"
)


// 获取用户测ip地址 
func GetClientIp(r *http.Request) string {
	if ips := r.Header.Get("X-Forwarded-For"); ips != "" {
		arr := strings.Split(ips, ",")
		return arr[0]
	}

	return ""
}