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

// ==========================================================================

var testCases = []apitests.Test{
	{
		Name:           "set count to 1",
		URL:            "/setProduct/demo@example.net/fake-01/1",
		Method:         "PUT",
		Body:           "",
		CheckBody:      `fake-01":1`,
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "set count to 34",
		URL:            "/setProduct/demo@example.net/fake-01/34",
		Method:         "PUT",
		Body:           "",
		CheckBody:      `fake-01":34`,
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "set count to zero",
		URL:            "/setProduct/demo@example.net/fake-99/0",
		Method:         "PUT",
		Body:           "",
		CheckBody:      `fake-99`,
		CheckBodyCount: 0,
		CheckStatus:    200,
	},
	{
		Name:           "set count to blah",
		URL:            "/setProduct/demo@example.net/fake-77/blah",
		Method:         "PUT",
		Body:           "",
		CheckBody:      "setProductCount failed",
		CheckBodyCount: 1,
		CheckStatus:    500,
	},
	{
		Name:           "set count to -1",
		URL:            "/setProduct/demo@example.net/fake-77/-1",
		Method:         "PUT",
		Body:           "",
		CheckBody:      "negative",
		CheckBodyCount: 1,
		CheckStatus:    500,
	},
	{
		Name:           "get cart",
		URL:            "/get/demo@example.net",
		Method:         "GET",
		Body:           "",
		CheckBody:      `"fake-01":34`,
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "get cart for invalid user",
		URL:            "/get/invalid@example.net",
		Method:         "GET",
		Body:           "",
		CheckBody:      "\\{\\}.*invalid@example.net",
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "submit cart",
		URL:            "/submit",
		Method:         "POST",
		Body:           `"demo@example.net"`,
		CheckBody:      `"status":"new"`,
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "clear cart",
		URL:            "/clear/demo@example.net",
		Method:         "PUT",
		Body:           "",
		CheckBody:      `"fake-01":34`,
		CheckBodyCount: 0,
		CheckStatus:    200,
	},
	{
		Name:           "submit empty cart",
		URL:            "/submit",
		Method:         "POST",
		Body:           `"demo@example.net"`,
		CheckBody:      ``,
		CheckBodyCount: 0,
		CheckStatus:    400,
	},
	{
		Name:           "submit cart for invalid user",
		URL:            "/submit",
		Method:         "POST",
		Body:           `"david_bowie@example.net"`,
		CheckBody:      ``,
		CheckBodyCount: 0,
		CheckStatus:    400,
	},
}
