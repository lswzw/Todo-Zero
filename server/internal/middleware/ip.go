package middleware

import (
	"net"
	"net/http"
	"strings"
)

var trustedProxyIPs []net.IP

func SetTrustedProxies(ips []string) {
	trustedProxyIPs = nil
	for _, ipStr := range ips {
		if ip := net.ParseIP(ipStr); ip != nil {
			trustedProxyIPs = append(trustedProxyIPs, ip)
		}
	}
}

func isTrustedProxy(ip string) bool {
	if len(trustedProxyIPs) == 0 {
		return false
	}
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	for _, trusted := range trustedProxyIPs {
		if parsedIP.Equal(trusted) {
			return true
		}
	}
	return false
}

func getClientIP(r *http.Request) string {
	if isTrustedProxy(r.RemoteAddr) {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			if len(parts) > 0 {
				ip := strings.TrimSpace(parts[0])
				if ip != "" {
					return ip
				}
			}
		}
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			return strings.TrimSpace(xri)
		}
	}

	addr := r.RemoteAddr
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}
