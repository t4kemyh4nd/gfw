package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/t4kemyh4nd/gfw/scanner"
)

const (
	proxyPort     = ":9001"
	defaultPort   = ":8000"
	defaultTarget = "127.0.0.1"
	scheme        = "http"
)

type Blacklist struct {
	headers   []string
	ips       []string
	locations []string
}

var Bl Blacklist

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		resp.Body = ioutil.NopCloser(bytes.NewReader(b))
	}

	//SET CUSTOM RESPONSE HEADERS HERE
	resp.Header.Set("Server", "gfw")

	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	return resp, nil
}

func main() {
	var origin, _ = url.Parse(scheme + "://" + defaultTarget + defaultPort)
	Bl.headers = []string{"X-FORBIDDEN"}
	Bl.locations = []string{"/forbidden"}
	Bl.ips = []string{"127.0.0.1"}

	var director = func(req *http.Request) {
		CleanHeaders(req)
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.Header.Add("X-Firewall", "1")
		req.URL.Scheme = scheme
		req.URL.Host = origin.Host
	}

	var reverseProxy = &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		CleanPath(r)
		if !IsForbiddenPath(r) {
			if scanner.ScanForXSS(r) && scanner.ScanForRCE(r) && scanner.ScanForSqli(r) && !IsIPBlacklisted(r) {
				reverseProxy.Transport = &transport{http.DefaultTransport}
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
