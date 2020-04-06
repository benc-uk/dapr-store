package main

//
// Basic REST API microservice, template/reference code
// Ben Coleman, July 2019, v1
//

import (
	"encoding/json"
	"net/http"
)

//
// Call Orders service
//
func (h API) ordersAPI(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "application/json")
	var example = make(map[string]string)
	example["message"] = "Hello from orders service"

	exampleJSON, err := json.Marshal(example)
	if err != nil {
		http.Error(resp, "Failure", http.StatusInternalServerError)
	}

	resp.Write(exampleJSON)
}
