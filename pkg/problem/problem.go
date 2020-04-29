// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// RFC-7807 implementation for sending standard format API errors
// ----------------------------------------------------------------------------

package problem

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Problem in RFC-7807 format
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

//
// New creates a RFC 7807 problem object
//
func New(url, title string, status int, detail, instance string) *Problem {
	return &Problem{url, title, status, detail, instance}
}

//
// HTTPSend sends a RFC 7807 problem object as HTTP response
//
func (p *Problem) Send(resp http.ResponseWriter) {
	log.Printf("### API %s", p.Error())
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(p.Status)
	json.NewEncoder(resp).Encode(p)
}

//
// NewAPIProblem creates a Problem based on either a HTTP resp or an error
//
func NewAPIProblem(url string, title string, instance string, apiResp *http.Response, err error) *Problem {
	var p *Problem
	if err != nil {
		p = New(url, title, 500, err.Error(), instance)
	} else if apiResp != nil {
		p = New(url, title, apiResp.StatusCode, http.StatusText(apiResp.StatusCode), instance)
	} else {
		p = New(url, title, 500, "Other error occurred", instance)
	}
	return p
}

// Implement error interface
func (p Problem) Error() string {
	return fmt.Sprintf("Problem: Type: '%s', Title: '%s', Status: '%d', Detail: '%s', Instance: '%s'", p.Type, p.Title, p.Status, p.Detail, p.Instance)
}
