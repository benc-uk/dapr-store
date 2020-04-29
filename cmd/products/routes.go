// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for products
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"net/http"

	"github.com/benc-uk/dapr-store/pkg/problem"
	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/get/{id}", api.getProduct)
	router.HandleFunc("/catalog", api.getCatalog)
	router.HandleFunc("/offers", api.getOffers)
	router.HandleFunc("/search/{query}", api.searchProducts)
}

//
// Return a single product
//
func (api API) getProduct(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	products, err := api.service.QueryProducts("ID", vars["id"])
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	// Handle no results
	if len(products) < 1 {
		prob := problem.New("err://products-db", "Not found", 404, "Product id: '"+vars["id"]+"' not found in DB", serviceName)
		prob.Send(resp)
		return
	}

	json, _ := json.Marshal(products)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(json)
}

//
// Return the product catalog
//
func (api API) getCatalog(resp http.ResponseWriter, req *http.Request) {
	products, err := api.service.AllProducts()
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	json, _ := json.Marshal(products)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(json)
}

//
// Return the products on offer
//
func (api API) getOffers(resp http.ResponseWriter, req *http.Request) {
	products, err := api.service.QueryProducts("onoffer", "1")
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	json, _ := json.Marshal(products)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(json)
}

//
// Search the products table
//
func (api API) searchProducts(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	products, err := api.service.SearchProducts(vars["query"])
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	json, _ := json.Marshal(products)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(json)
}
