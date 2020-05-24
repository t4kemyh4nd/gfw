package scanner

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Scanner interface {
	getGETValues() map[string][]string
	removeGETMalChars() bool
}

//PLEASE USE PARAMETERIZED SQL QUERIES IN YOUR CODE
type SQLiscanner struct {
	req *http.Request
}

func (s SQLiscanner) getGETValues() map[string][]string {
	var queryMap = make(map[string][]string)

	if s.req.Method == "GET" {
		queryMap = s.req.URL.Query()
	}

	if s.req.Method == "POST" {
		queryMap = s.req.URL.Query()
	}

	return queryMap
}

func (s SQLiscanner) removeGETMalChars() bool {
	var malChars = []string{"'", "--", "\"", "||"}
	//sqli regex
	var sqlRegex = regexp.MustCompile(`(\bunion(\(*|\s{1,})select\s{1,}.*\s{1,}from(\(*|\s{1,})|\binsert\s{1,}into\s{1,}\({0,1}.*\){0,1}\s{1,}values\s*\({0,1}|(#|--)$)`)
	var flag bool = true

	var queryMap = s.getGETValues()

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

func ScanForSqli(req *http.Request) {
	var sqlsicanner Scanner
	sqlsicanner = SQLiscanner{req}

	sqlsicanner.removeGETMalChars()
}
