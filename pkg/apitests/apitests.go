// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Helper for running API/route based tests
// ----------------------------------------------------------------------------

package apitests

import (
	"io/ioutil"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type Test struct {
	Name           string
	URL            string
	Method         string
	Body           string
	CheckBody      string
	CheckBodyCount int
	CheckStatus    int
}

func Run(t *testing.T, router *mux.Router, testCases []Test) {
	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			req := httptest.NewRequest(test.Method, test.URL, strings.NewReader(test.Body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Length", strconv.Itoa(len(test.Body)))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			body, _ := ioutil.ReadAll(rec.Result().Body)
			if rec.Result().StatusCode != test.CheckStatus {
				t.Errorf("Got status %d wanted %d", rec.Result().StatusCode, test.CheckStatus)
				return
			}

			if test.CheckBody != "" {
				bodyCheckRegex := regexp.MustCompile(test.CheckBody)
				matches := bodyCheckRegex.FindAllStringIndex(string(body), -1)

				if len(matches) != test.CheckBodyCount {
					t.Errorf("'%s' not found %d times in body ", test.CheckBody, test.CheckBodyCount)
					t.Logf(" BODY: %s", body)
					return
				}
			}
			//t.Logf("'%s' passed.", test.Name)
		})
	}
}
