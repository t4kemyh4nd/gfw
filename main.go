package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gfw/scanner"
)

const (
	proxyPort     = ":9001"
	defaultPort   = ":8000"
	defaultTarget = "127.0.0.1"
)

type Blacklist struct {
	headers   []string
	ips       []string
	locations []string
}

var Bl Blacklist

func main() {
	var origin, _ = url.Parse("http://127.0.0.1:8000")
	Bl.headers = []string{"X-FORBIDDEN"}
	Bl.locations = []string{"/forbidden"}

	var director = func(req *http.Request) {
		CleanHeaders(req)
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.Header.Add("X-FIREWALL", "1")
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	var reverseProxy = &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		CleanPath(r)
		if !IsForbiddenPath(r) {
			if scanner.ScanForXSS(r) && scanner.ScanForRCE(r) && scanner.ScanForSqli(r) {
				reverseProxy.ServeHTTP(w, r)
			} else {
				w.WriteHeader(403)
			}
		} else {
			w.WriteHeader(403)
		}
	})

	log.Fatal(http.ListenAndServe(proxyPort, nil))
}
