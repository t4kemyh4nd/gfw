package scanner

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Scanner interface {
	getGETValues(*http.Request) map[string][]string
	removeGETMalChars(*http.Request) bool
}

//PLEASE USE PARAMETERIZED SQL QUERIES IN YOUR CODE
type SQLiscanner struct {
	req *http.Request
}

func (s SQLiscanner) getGETValues(req *http.Request) map[string][]string {
	var queryMap = make(map[string][]string)

	if req.Method == "GET" {
		queryMap = req.URL.Query()
	}

	if req.Method == "POST" {
		queryMap = req.URL.Query()
	}

	return queryMap
}

func (s SQLiscanner) removeGETMalChars(req *http.Request) bool {
	var malChars = []string{"'", "--", "\"", "||"}
	//sqli regex
	var sqlRegex = regexp.MustCompile(`(\bunion(\(*|\s{1,})select\s{1,}.*\s{1,}from(\(*|\s{1,})|\binsert\s{1,}into\s{1,}\({0,1}.*\){0,1}\s{1,}values\s*\({0,1}|(#|--)$)`)
	var flag bool = true

	var queryMap = s.getGETValues(req)

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

	req.URL.RawQuery = getQuery.Encode()
	return flag
}

func ScanForSqli(req *http.Request) {
	var sqlsicanner Scanner
	sqlsicanner = SQLiscanner{req}

	sqlsicanner.removeGETMalChars(req)
}
