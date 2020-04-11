// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/benc-uk/dapr-store/common"

	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/get/{id}", api.getProduct)
	router.HandleFunc("/catalog", api.getCatalog)
	router.HandleFunc("/offers", api.getOffers)
	router.HandleFunc("/reload", api.reloadProducts)
}

//
// Return a single product
//
func (api API) getProduct(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	productJSON := []byte("{}")
	if _, exists := products[id]; exists {
		productJSON, _ = json.Marshal(products[id])
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(productJSON)
}

//
// Return the product catalog
//
func (api API) getCatalog(resp http.ResponseWriter, req *http.Request) {
	productsJSON, _ := json.Marshal(products)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(productsJSON)
}

//
// Return the products on offer
//
func (api API) getOffers(resp http.ResponseWriter, req *http.Request) {
	offers := make(map[string]common.Product)

	for id, prod := range products {
		if prod.OnOffer == true {
			offers[id] = prod
		}
	}

	offersJSON, _ := json.Marshal(offers)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(offersJSON)
}

//
// Return the products on offer
//
func (api API) reloadProducts(resp http.ResponseWriter, req *http.Request) {
	// Load index which is just an array of keys/ids
	data, err := common.GetState(resp, daprPort, daprStateStore, serviceName, "index")
	if err != nil {
		return
	}

	productIds := []int{}
	err = json.Unmarshal(data, &productIds)
	if err != nil {
		common.Problem{"json-error", "Error loading product index", 500, err.Error(), serviceName}.HttpSend(resp)
		return
	}

	// Load each object, and push into array
	products = make(map[string]common.Product)
	for _, id := range productIds {
		p := common.Product{}
		data, err := common.GetState(resp, daprPort, daprStateStore, serviceName, strconv.Itoa(id))
		if err != nil {
			log.Printf("### Error loading product '%v' %+v\n", id, err.Error())
			continue
		}
		err = json.Unmarshal(data, &p)
		if err != nil {
			log.Printf("### Error decoding product '%v' %+v\n", id, err.Error())
			continue
		}
		products[strconv.Itoa(id)] = p
	}

	//log.Printf("%+v", products)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(`{"message": "done", "count": "blah"}`))
}
