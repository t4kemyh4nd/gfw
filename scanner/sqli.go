package scanner

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Scanner interface {
	getGETValues(*http.Request) map[string][]string
	removeGETMalChars(*http.Request) bool
}

//PLEASE USE PARAMETERIZED SQL QUERIES IN YOUR CODE
type SQLiScanner struct {
	req *http.Request
}

func (s SQLiScanner) getGETValues(req *http.Request) map[string][]string {
	var queryMap = make(map[string][]string)

	if req.Method == "GET" {
		queryMap = req.URL.Query()
	}

	if req.Method == "POST" {
		queryMap = req.URL.Query()
	}

	return queryMap
}

func (s SQLiScanner) removeGETMalChars(req *http.Request) bool {
	var malChars = []string{"'", "--", "\"", ";"}
	var queryMap = s.getGETValues(req)

	for keys, values := range queryMap {
		for index, param := range values {
			fmt.Println("params:", param)
			for strings.Contains(param, "/*") || strings.Contains(param, "*/") || strings.Contains(param, "#") {
				param = strings.ReplaceAll(param, "/*", "")
				param = strings.ReplaceAll(param, "*/", " ")
				param = strings.ReplaceAll(param, "#", " ")
			}
			for _, m := range malChars {
				if strings.Contains(param, m) {
					param = strings.ReplaceAll(param, m, "\\"+m)
				}
			}
			values[index] = param
		}
		queryMap[keys] = values
	}

	var getQuery = url.Values{}
	getQuery = queryMap

	req.URL.RawQuery = getQuery.Encode()
	return true
}

func ScanForSqli(req *http.Request) {
	var sqlsicanner Scanner
	sqlsicanner = SQLiScanner{req}

	sqlsicanner.removeGETMalChars(req)
}
