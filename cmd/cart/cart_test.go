// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Main set of tests for cart service and API
// ----------------------------------------------------------------------------

package main

import (
	"io"
	"log"
	"testing"

	"github.com/benc-uk/dapr-store/cmd/cart/mock"
	"github.com/benc-uk/go-rest-api/pkg/api"
	"github.com/benc-uk/go-rest-api/pkg/auth"
	"github.com/benc-uk/go-rest-api/pkg/httptester"
	"github.com/go-chi/chi/v5"
)

func TestCart(t *testing.T) {
	// Comment out to see logs
	log.SetOutput(io.Discard)

	// Mock of CartService
	mockCartSvc := &mock.CartService{}

	router := chi.NewRouter()
	api := API{
		api.NewBase("cart", "ignore", "ignore", true),
		mockCartSvc,
	}
	api.addRoutes(router, auth.NewPassthroughValidator())

	httptester.Run(t, router, testCases)
}

// ==========================================================================

var testCases = []httptester.TestCase{
	{
		Name:           "set count to 1",
		URL:            "/setProduct/mock@example.net/fake-01/1",
		Method:         "PUT",
		Body:           "",
		CheckBody:      `fake-01":1`,
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "set count to 34",
		URL:            "/setProduct/mock@example.net/fake-01/34",
		Method:         "PUT",
		Body:           "",
		CheckBody:      `fake-01":34`,
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "set count to zero",
		URL:            "/setProduct/mock@example.net/fake-99/0",
		Method:         "PUT",
		Body:           "",
		CheckBody:      `fake-99`,
		CheckBodyCount: 0,
		CheckStatus:    200,
	},
	{
		Name:           "set count to blah",
		URL:            "/setProduct/mock@example.net/fake-77/blah",
		Method:         "PUT",
		Body:           "",
		CheckBody:      "invalid syntax",
		CheckBodyCount: 1,
		CheckStatus:    400,
	},
	{
		Name:           "set count to -1",
		URL:            "/setProduct/mock@example.net/fake-77/-1",
		Method:         "PUT",
		Body:           "",
		CheckBody:      "product count",
		CheckBodyCount: 1,
		CheckStatus:    500,
	},
	// {
	// 	Name:           "set count for non-existing user",
	// 	URL:            "/setProduct/dontexist@example.net/bloop88/1",
	// 	Method:         "PUT",
	// 	Body:           "",
	// 	CheckBody:      ``,
	// 	CheckBodyCount: 0,
	// 	CheckStatus:    200,
	// },
	{
		Name:           "get cart",
		URL:            "/get/mock@example.net",
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
		Body:           `"mock@example.net"`,
		CheckBody:      `"status":"new"`,
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "clear cart",
		URL:            "/clear/mock@example.net",
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
		Body:           `"mock@example.net"`,
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
