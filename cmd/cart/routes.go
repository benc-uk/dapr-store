// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for cart
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/benc-uk/dapr-store/pkg/auth"
	"github.com/benc-uk/dapr-store/pkg/problem"

	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/setProduct/{username}/{productId}/{count}", auth.JWTValidator(api.setProductCount)).Methods("PUT")
	router.HandleFunc("/get/{username}", auth.JWTValidator(api.getCart)).Methods("GET")
	router.HandleFunc("/submit", auth.JWTValidator(api.submitCart)).Methods("POST")
	router.HandleFunc("/clear/{username}", auth.JWTValidator(api.clearCart)).Methods("PUT")
}

//
//
//
func (api API) setProductCount(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cart, _ := api.service.Get(vars["username"])

	count, err := strconv.Atoi(vars["count"])
	if err != nil {
		problem.New("err://invalid-count", "setProductCount failed", 500, err.Error(), api.ServiceName).Send(resp)
		return
	}

	err = api.service.SetProductCount(cart, vars["productId"], count)
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(cart)
	log.Printf("cart %s", json)
	_, _ = resp.Write(json)
}

//
//
//
func (api API) getCart(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cart, err := api.service.Get(vars["username"])
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(cart)
	log.Printf("cart %s", json)
	_, _ = resp.Write(json)
}

//
//
//
func (api API) clearCart(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cart, err := api.service.Get(vars["username"])
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	err = api.service.Clear(cart)
	if err != nil {
		log.Printf("### Warning failed to clear cart %s", err)
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(cart)
	log.Printf("cart %s", json)
	_, _ = resp.Write(json)
}

//
//
//
func (api API) submitCart(resp http.ResponseWriter, req *http.Request) {
	cl, _ := strconv.Atoi(req.Header.Get("content-length"))
	if cl <= 0 {
		problem.New("err://body-missing", "Zero length body", 400, "Zero length body", api.ServiceName).Send(resp)
		return
	}

	username := ""
	err := json.NewDecoder(req.Body).Decode(&username)

	// Some basic validation and checking on what we've been posted
	if err != nil {
		problem.New("err://json-decode", "Malformed JSON", 400, "JSON could not be decoded", api.ServiceName).Send(resp)
		return
	}
	if username == "" {
		problem.New("err://json-error", "Malformed JSON", 400, "Post should include username", api.ServiceName).Send(resp)
		return
	}

	cart, err := api.service.Get(username)
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	order, err := api.service.Submit(*cart)
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	// Send the _order_ back, created from submitting the cart
	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(order)
	_, _ = resp.Write(json)
}
