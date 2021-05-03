package scanner

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Scanner interface {
	getValues() map[string][]string
	removeMalChars() bool
}

//PLEASE USE PARAMETERIZED SQL QUERIES IN YOUR CODE
type SQLiscanner struct {
	req *http.Request
}

func (s SQLiscanner) getValues() map[string][]string {
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

func (s SQLiscanner) removeMalChars() bool {
	var malChars = []string{"'", "--", "\"", "||"}
	//sqli regex
	var sqlRegex = regexp.MustCompile(`(\bunion(\(*|\s{1,})select\s{1,}(.*|from(\(*|\s{1,}).*)|\binsert\s{1,}into\s{1,}\({0,1}.*\){0,1}\s{1,}values\s*\({0,1}|)(#|--)$`)
	var flag bool = true

	var queryMap = s.getValues()

	for keys, values := range queryMap {
		for index, param := range values {
			for strings.Contains(param, "/*") || strings.Contains(param, "*/") || strings.Contains(param, "#") {
				param = strings.ReplaceAll(param, "/*", "")
				param = strings.ReplaceAll(param, "*/", " ")
			}
			for _, m := range malChars {
				if strings.Contains(param, m) {
					param = strings.ReplaceAll(strings.ToLower(param), m, "\\"+m)
				}

			}
			flag = !sqlRegex.MatchString(strings.ToLower(param))
			values[index] = param
		}
		queryMap[keys] = values
	}

	var getQuery = url.Values{}
	getQuery = queryMap

	s.req.URL.RawQuery = getQuery.Encode()
	return flag
}

func ScanForSqli(req *http.Request) bool {
	var sqlsicanner Scanner = &SQLiscanner{req}

	return sqlsicanner.removeMalChars()
}
