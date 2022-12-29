// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for products
// ----------------------------------------------------------------------------

package main

import (
	"errors"
	"net/http"

	"github.com/benc-uk/go-rest-api/pkg/problem"
	"github.com/go-chi/chi/v5"
)

// All routes we need should be registered here
func (api API) addRoutes(router chi.Router) {
	router.Get("/get/{id}", api.getProduct)
	router.Get("/catalog", api.getCatalog)
	router.Get("/offers", api.getOffers)
	router.Get("/search/{query}", api.searchProducts)
}

// Return a single product
func (api API) getProduct(resp http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	products, err := api.service.QueryProducts("ID", id)
	if err != nil {
		problem.Wrap(500, req.RequestURI, id, err).Send(resp)

		return
	}

	// Handle no results
	if len(products) < 1 {
		problem.Wrap(404, req.RequestURI, id, errors.New("product not found")).Send(resp)

		return
	}

	api.ReturnJSON(resp, products[0])
}

// Return the product catalog
func (api API) getCatalog(resp http.ResponseWriter, req *http.Request) {
	products, err := api.service.AllProducts()
	if err != nil {
		problem.Wrap(500, req.RequestURI, "catalog", err).Send(resp)

		return
	}

	api.ReturnJSON(resp, products)
}

// Return the products on offer
func (api API) getOffers(resp http.ResponseWriter, req *http.Request) {
	products, err := api.service.QueryProducts("onoffer", "1")
	if err != nil {
		problem.Wrap(500, req.RequestURI, "offers", err).Send(resp)

		return
	}

	api.ReturnJSON(resp, products)
}

// Search the products table
func (api API) searchProducts(resp http.ResponseWriter, req *http.Request) {
	query := chi.URLParam(req, "query")

	products, err := api.service.SearchProducts(query)
	if err != nil {
		problem.Wrap(500, req.RequestURI, query, err).Send(resp)

		return
	}

	api.ReturnJSON(resp, products)
}
