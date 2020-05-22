package main

import (
	"fmt"
	"net/http"
)

type SQLiScanner struct {
	req *http.Request
}

func (s SQLiScanner) getValues(req *http.Request) []string {

}

func (s SQLiScanner) getMalChars(req *http.Request) []string {

}

func ScanForSqli(req *http.Request) {
	var sqlsicanner Scanner
	sqlsicanner = &SQLiScanner{req}

	var queryMap = make(map[string][]string)

	if req.Method == "GET" {
		queryMap = req.URL.Query()
	}

	if req.Method == "POST" {
		//scan POST parameters
	}
	fmt.Println(queryMap)
}
