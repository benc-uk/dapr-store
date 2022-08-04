// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for orders
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"net/http"

	"github.com/benc-uk/dapr-store/pkg/auth"
	"github.com/benc-uk/dapr-store/pkg/problem"
	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/get/{id}", auth.JWTValidator(api.getOrder))
	router.HandleFunc("/getForUser/{username}", auth.JWTValidator(api.getOrdersForUser))
}

//
// Fetch existing order by id
//
func (api API) getOrder(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	order, err := api.service.GetOrder(vars["id"])
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)

		return
	}

	resp.Header().Set("Content-Type", "application/json")

	json, _ := json.Marshal(order)
	_, _ = resp.Write(json)
}

//
// Fetch all orders for a given user
//
func (api API) getOrdersForUser(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	orders, err := api.service.GetOrdersForUser(vars["username"])
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)

		return
	}

	resp.Header().Set("Content-Type", "application/json")

	json, _ := json.Marshal(orders)
	_, _ = resp.Write(json)
}
