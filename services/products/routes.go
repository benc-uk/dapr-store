// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service
// ----------------------------------------------------------------------------

package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/benc-uk/dapr-store/common"
	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/util/json"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/get/{id}", api.getProduct)
	router.HandleFunc("/catalog", api.getCatalog)
	router.HandleFunc("/offers", api.getOffers)
}

//
// Return a single product
//
func (api API) getProduct(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	p := common.Product{}
	stmt, _ := db.Prepare("SELECT * FROM products WHERE ID = ?")
	defer stmt.Close()
	err := stmt.QueryRow(vars["id"]).Scan(&p.ID, &p.Name, &p.Description, &p.Cost, &p.Image, &p.OnOffer)
	if err == sql.ErrNoRows {
		common.Problem{"products-db", "Product " + vars["id"] + " not found in DB", 404, err.Error(), serviceName}.HttpSend(resp)
		return
	}
	if err != nil {
		log.Printf("### Products DB error: %+v\n", err)
		common.Problem{"products-db", "Products DB error", 500, err.Error(), serviceName}.HttpSend(resp)
		return
	}

	productJSON, _ := json.Marshal(p)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(productJSON)
}

//
// Return the product catalog
//
func (api API) getCatalog(resp http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		common.Problem{"database", "Error querying products", 500, err.Error(), serviceName}.HttpSend(resp)
		return
	}

	returnProducts(rows, resp)
}

//
// Return the products on offer
//
func (api API) getOffers(resp http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("SELECT * FROM products WHERE onoffer = true")
	if err != nil {
		common.Problem{"database", "Error querying products", 500, err.Error(), serviceName}.HttpSend(resp)
		return
	}

	returnProducts(rows, resp)
}

//
//
//
func returnProducts(rows *sql.Rows, resp http.ResponseWriter) {
	products := []common.Product{}
	defer rows.Close()
	for rows.Next() {
		p := common.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Cost, &p.Image, &p.OnOffer)
		if err != nil {
			common.Problem{"database", "Error reading row", 500, err.Error(), serviceName}.HttpSend(resp)
			return
		}
		products = append(products, p)
	}

	productsJSON, _ := json.Marshal(products)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(productsJSON)
}
