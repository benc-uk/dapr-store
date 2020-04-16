// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service
// ----------------------------------------------------------------------------

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// All routes we need should be here
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/get/{id}", api.getExample)
}

//
// Example
//
func (api API) getExample(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(fmt.Sprintf(`{"message": "echo, id was %v"}`, vars["id"])))
}
