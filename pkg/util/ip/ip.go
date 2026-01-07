package ip

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP 获取客户端真实 IP
func GetClientIP(r *http.Request) string {
	// 1. 尝试从 X-Forwarded-For 获取
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For 可能包含多个 IP，第一个是真实客户端 IP
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			if clientIP != "" {
				return clientIP
			}
		}
	}

	// 2. 尝试从 X-Real-IP 获取
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// 3. 回退到 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// 如果 RemoteAddr 不包含端口（虽然不太可能），直接返回
		return r.RemoteAddr
	}
	return ip
}
