// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// RFC-7807 implementation for sending standard format API errors
// ----------------------------------------------------------------------------

package problem

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

// HTTP400 is a fake static response for returning 400s to client
var HTTP400 = &http.Response{StatusCode: 400, Status: "Bad request"}

// HTTP404 is a fake static response for returning 404s to client
var HTTP404 = &http.Response{StatusCode: 404, Status: "Not found"}

//
// HTTPSend sends a RFC 7807 problem object as HTTP response
//
func (p Problem) HTTPSend(resp http.ResponseWriter) {
	log.Printf("### Sending API problem to client: %+v", p)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(p.Status)
	json.NewEncoder(resp).Encode(p)
}

//
// Send creates a standardized format problem from either the 'err' or 'apiResp' and sends it
//
func Send(title string, url string, resp http.ResponseWriter, apiResp *http.Response, err error, instance string) {
	if err != nil {
		p := Problem{url, title, 500, err.Error(), instance}
		p.HTTPSend(resp)
	} else if apiResp != nil {
		p := Problem{url, title, apiResp.StatusCode, http.StatusText(apiResp.StatusCode), instance}
		p.HTTPSend(resp)
	} else {
		p := Problem{url, title, 500, "Other error occurred", instance}
		p.HTTPSend(resp)
	}
}
