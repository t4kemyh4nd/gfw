package main

import (
	"net"
	"net/http"
	"strings"
)

func CleanHeaders(req *http.Request) {
	//REMOVE REQUEST DESYNC HEADERS
	if req.Header.Get("Content-Length") != "" && req.Header.Get("Transfer-Encoding") != "" {
		req.Header.Del("Content-Length")
	}
}

func CleanPath(req *http.Request) {
	var path = req.URL.Path

	//REMOVE DOUBLE SLASHES AND PATH TRAVERSALS
	for strings.Contains(path, "../") || strings.Contains(path, "//") || strings.Contains(path, "..;/") {
		path = strings.ReplaceAll(path, "//", "/")
		path = strings.ReplaceAll(path, "../", "/")
	}

	req.URL.Path = path
}

func IsForbiddenPath(req *http.Request) bool {
	var flag bool
	for _, path := range Bl.locations {
		if req.URL.Path == path {
			flag = true
		} else {
			flag = false
		}
	}
	return flag
}

func IsIPBlacklisted(req *http.Request) bool {
	if len(Bl.ips) > 0 {
		var ip string
		var flag bool = false

		if req.Header.Get("X-Forwarded-For") != "" {
			ip = req.Header.Get("X-Forwarded-For")
		} else if req.Header.Get("X-Forwarded-IP") != "" {
			ip = req.Header.Get("X-Forwarded-IP")
		} else {
			ip, _, _ = net.SplitHostPort(req.RemoteAddr)
		}

		reqIP := net.ParseIP(ip)

		for _, ip := range Bl.ips {
			if reqIP.String() == net.ParseIP(ip).String() {
				flag = true
				break
			}
		}

		return flag
	}

	return false
}
