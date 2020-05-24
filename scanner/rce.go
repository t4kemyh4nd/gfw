package scanner

import (
	"net/http"
	"regexp"
	"strings"
)

type RCEscanner struct {
	req *http.Request
}

func (s RCEscanner) getGETValues() map[string][]string {
	var queryMap = make(map[string][]string)

	if s.req.Method == "GET" {
		queryMap = s.req.URL.Query()
	}

	if s.req.Method == "POST" {
		queryMap = s.req.URL.Query()
	}

	return queryMap
}

func (s RCEscanner) removeGETMalChars() bool {
	var rceRegex = regexp.MustCompile(`(;|&&|\|\|)(\s{0,}|\{IFS\})(sleep|curl|wget|netcat|nc|nslookup|ping|cat|touch)(\s{0,}|\{IFS\}).*(;|&&|\|\|)`)
	var queryMap = s.getGETValues()
	var flag bool = true

	for _, values := range queryMap {
		for _, param := range values {
			flag = !rceRegex.MatchString(strings.ToLower(param))
		}
	}

	return flag
}

func ScanForRCE(req *http.Request) bool {
	var rcescanner Scanner
	rcescanner = &RCEscanner{req}

	return rcescanner.removeGETMalChars()
}
