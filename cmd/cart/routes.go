// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for cart
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/dapr-store/pkg/auth"
	"github.com/benc-uk/dapr-store/pkg/problem"

	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/submit", auth.AuthMiddleware(api.submitOrder)).Methods("POST")
}

//
// Create a new order
//
func (api API) submitOrder(resp http.ResponseWriter, req *http.Request) {
	cl, _ := strconv.Atoi(req.Header.Get("content-length"))
	if cl <= 0 {
		problem.New("err://body-missing", "Zero length body", 400, "Zero length body", api.ServiceName).Send(resp)
		return
	}

	order := spec.Order{}
	err := json.NewDecoder(req.Body).Decode(&order)

	// Some basic validation and checking on what we've been posted
	if err != nil {
		problem.New("err://json-decode", "Malformed order JSON", 400, "JSON could not be decoded", api.ServiceName).Send(resp)
		return
	}
	if err := spec.Validate(order); err != nil {
		problem.New("err://json-error", "Malformed order JSON", 400, "Order failed validation, check spec", api.ServiceName).Send(resp)
		return
	}

	api.service.SubmitOrder(&order)

	// Send the order back, but this time it will have an id
	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(order)
	resp.Write(json)
}
