package scanner

import (
	"fmt"
	"net/http"
	"regexp"
)

type RCEscanner struct {
	req *http.Request
}

func (s RCEscanner) getGETValues(req *http.Request) map[string][]string {
	var queryMap = make(map[string][]string)

	if req.Method == "GET" {
		queryMap = req.URL.Query()
	}

	if req.Method == "POST" {
		queryMap = req.URL.Query()
	}

	return queryMap
}

func (s RCEscanner) removeGETMalChars(req *http.Request) bool {
	var rceRegex = regexp.MustCompile(`(;|&&|\|\|)(\s{0,}|\{IFS\})(sleep|curl|wget|netcat|nc|nslookup|ping|cat|touch)(\s{0,}|\{IFS\}).*(;|&&|\|\|)`)
	var queryMap = s.getGETValues(req)
	var flag bool = true

	for _, values := range queryMap {
		for _, param := range values {
			flag = !rceRegex.MatchString(param)
			fmt.Println(param)
			fmt.Println(flag)
		}
	}

	return flag
}

func ScanForRCE(req *http.Request) bool {
	var rcescanner Scanner
	rcescanner = RCEscanner{req}

	return rcescanner.removeGETMalChars(req)
}
