package scanner

import (
	"net/http"
	"regexp"
	"strings"
)

type RCEscanner struct {
	req *http.Request
}

func (s RCEscanner) getValues() map[string][]string {
	if s.req.Method == "GET" {
		return s.req.URL.Query()
	} else {
		if err := s.req.ParseForm(); err != nil {
			panic("Couldn't parse form")
		} else {
			return s.req.PostForm
		}
	}
}

func (s RCEscanner) removeMalChars() bool {
	var rceRegex = regexp.MustCompile(`(;|&&|\|\|)(\s{0,}|\{IFS\})(sleep|curl|wget|netcat|nc|nslookup|ping|cat|touch)(\s{0,}|\{IFS\}).*(;|&&|\|\|)`)
	var queryMap = s.getValues()
	var flag bool = true

	for _, values := range queryMap {
		for _, param := range values {
			flag = !rceRegex.MatchString(strings.ToLower(param))
		}
	}

	return flag
}

func ScanForRCE(req *http.Request) bool {
	var rcescanner Scanner = &RCEscanner{req}

	return rcescanner.removeMalChars()
}
