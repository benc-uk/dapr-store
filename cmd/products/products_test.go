// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Main set of tests for products service and API
// ----------------------------------------------------------------------------

package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/benc-uk/dapr-store/cmd/products/mock"
	"github.com/benc-uk/dapr-store/pkg/api"
	"github.com/benc-uk/dapr-store/pkg/apitests"
	"github.com/gorilla/mux"
)

func TestProducts(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	// Mock of ProductsService
	mockProductSvc := &mock.ProductsService{}

	router := mux.NewRouter()
	api := API{
		api.NewBase("products", "ignore", "ignore", true, router),
		mockProductSvc,
	}
	api.addRoutes(router)

	apitests.Run(t, router, testCases)
}

var testCases = []apitests.Test{
	{
		Name:           "search for 'Hat'",
		URL:            "/search/Hat",
		Method:         "GET",
		Body:           "",
		CheckBody:      `prd1`,
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "search for cheese",
		URL:            "/search/cheese",
		Method:         "GET",
		Body:           "",
		CheckBody:      "\\[\\]",
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "get product prd3",
		URL:            "/get/prd3",
		Method:         "GET",
		Body:           "",
		CheckBody:      "prd3",
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "get on offer products",
		URL:            "/offers",
		Method:         "GET",
		Body:           "",
		CheckBody:      "prd2",
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "get all products",
		URL:            "/catalog",
		Method:         "GET",
		Body:           `"id"`,
		CheckBodyCount: 3,
		CheckStatus:    200,
	},
	{
		Name:           "get non-existent product",
		URL:            "/get/999",
		Method:         "GET",
		Body:           "",
		CheckBody:      "",
		CheckBodyCount: 0,
		CheckStatus:    404,
	},
}
