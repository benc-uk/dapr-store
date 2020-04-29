// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for products
// ----------------------------------------------------------------------------

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/benc-uk/dapr-store/pkg/models"
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

	rows, err := db.Query("SELECT * FROM products WHERE ID = ? LIMIT 1", vars["id"])
	if err != nil {
		problem.New("err://products-db", "Database query error", 500, err.Error(), serviceName).Send(resp)
		return
	}
	defer rows.Close()
	hasRow := rows.Next()
	if !hasRow {
		problem.New("err://products-db", "Not found", 404, "Product id: '"+vars["id"]+"' not found in DB", serviceName).Send(resp)
		return
	}

	p := models.Product{}
	rows.Scan(&p.ID, &p.Name, &p.Description, &p.Cost, &p.Image, &p.OnOffer)
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
		problem.New("err://products-db", "Error querying products", 500, err.Error(), serviceName).Send(resp)
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
		problem.New("err://products-db", "Error querying products", 500, err.Error(), serviceName).Send(resp)
		return
	}

	returnProducts(rows, resp)
}

//
// Search the products table
//
func (api API) searchProducts(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	rows, err := db.Query("SELECT * FROM products WHERE (description LIKE ? OR name LIKE ?)", "%"+vars["query"]+"%", "%"+vars["query"]+"%")
	if err != nil {
		problem.New("err://products-db", "Error querying products", 500, err.Error(), serviceName).Send(resp)
		return
	}

	returnProducts(rows, resp)
}

//
//
//
func returnProducts(rows *sql.Rows, resp http.ResponseWriter) {
	products := []models.Product{}
	defer rows.Close()
	for rows.Next() {
		p := models.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Cost, &p.Image, &p.OnOffer)
		if err != nil {
			problem.New("err://products-db", "Error reading row", 500, err.Error(), serviceName).Send(resp)
			return
		}
		products = append(products, p)
	}

	productsJSON, _ := json.Marshal(products)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(productsJSON)
}
