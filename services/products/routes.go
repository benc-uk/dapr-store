// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service
// ----------------------------------------------------------------------------

package main

import (
	"net/http"

	"github.com/benc-uk/dapr-store/common"
	"github.com/gorilla/mux"
)

// All routes we need should be here
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/get/{id}", api.getProduct)
}

//
// Return a single product
//
func (api API) getProduct(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	data, err := common.GetState(resp, daprPort, daprStateStore, serviceName, vars["id"])

	if err != nil {
		return // Error will have already been written to resp
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
}
