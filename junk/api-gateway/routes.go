package main

//
// Basic REST API microservice, template/reference code
// Ben Coleman, July 2019, v1
//

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/benc-uk/dapr-store/common"

	"github.com/gorilla/mux"
)

//
// Make a generic service to service call with Dapr
//
func (h API) daprProxy(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// Proxy call via Dapr sidecar
	// In format invoke/{service}/method/{methodPath}
	invokeURL := fmt.Sprintf(
		"http://localhost:%d/v1.0/invoke/%s/method/%s",
		daprPort,
		vars["service"],
		vars["restOfURL"],
	)
	fmt.Println("### Making proxy invoke call to:", invokeURL)
	proxyResp, err := http.Get(invokeURL)

	// Check major / network errors
	if err != nil {
		common.SendProblem(resp, "", "Dapr network error", 502, err.Error())
		return
	}

	// HTTP errors could be downstream is dead/404 or a error returned
	if proxyResp.StatusCode < 200 || proxyResp.StatusCode > 299 {
		common.SendProblem(resp, "", "Downstream service HTTP code not OK", proxyResp.StatusCode, "")
		return
	}

	// Read the body, we assume it's JSON
	defer proxyResp.Body.Close()
	proxyBody, err := ioutil.ReadAll(proxyResp.Body)

	// Yet more chance for a error
	if err != nil {
		common.SendProblem(resp, "", "Body IO error", 502, err.Error())
		return
	}

	// Done, mimic the content type, status and body
	resp.Header().Set("Content-Type", proxyResp.Header.Get("Content-Type"))
	resp.WriteHeader(proxyResp.StatusCode)
	resp.Write(proxyBody)
}
