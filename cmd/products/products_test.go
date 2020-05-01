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
