package main

import (
	"net/http"
	"strings"
)

func CleanHeaders(req *http.Request) {
	//REMOVE REQUEST DESYNC HEADERS
	if req.Header.Get("Content-Length") != "" && req.Header.Get("Transfer-Encoding") != "" {
		req.Header.Del("Transfer-Encoding")
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

/*
func BlackListIPaddress(req *http.Request bool) {
	return true
}
*/
