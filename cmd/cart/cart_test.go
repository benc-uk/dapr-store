// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Main set of tests for cart service and API
// ----------------------------------------------------------------------------

package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/benc-uk/dapr-store/cmd/cart/mock"
	"github.com/benc-uk/dapr-store/pkg/api"
	"github.com/benc-uk/dapr-store/pkg/apitests"
	"github.com/gorilla/mux"
)

func TestCart(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	// Mock of CartService
	mockCartSvc := &mock.CartService{}

	router := mux.NewRouter()
	api := API{
		api.NewBase("cart", "ignore", "ignore", true, router),
		mockCartSvc,
	}
	api.addRoutes(router)

	apitests.Run(t, router, testCases)
}
