// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// RFC-7807 implementation for sending standard format API errors
// ----------------------------------------------------------------------------

package common

import (
	"encoding/json"
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
// HttpSend sends a RFC 7807 problem object as HTTP response
//
func (p Problem) HttpSend(resp http.ResponseWriter) {
	log.Printf("### Sending API problem to client: %+v", p)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(p.Status)
	json.NewEncoder(resp).Encode(p)
}

//
// SendDaprProblem simple check for err (serious problem) or HTTP status code error
//
func SendDaprProblem(url string, resp http.ResponseWriter, daprResp *http.Response, err error, instance string) {
	if err != nil {
		p := Problem{url, "Dapr network error", 500, err.Error(), instance}
		p.HttpSend(resp)
	} else {
		p := Problem{url, "Dapr service error", daprResp.StatusCode, http.StatusText(daprResp.StatusCode), instance}
		p.HttpSend(resp)
	}
}
