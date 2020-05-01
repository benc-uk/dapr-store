// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Test cases for testing products API, run by products_test.go
// ----------------------------------------------------------------------------

package main

import "github.com/benc-uk/dapr-store/pkg/apitests"

var testCases = []apitests.Test{
	{
		"search for hats",
		"/search/hat",
		"GET",
		``,
		"\"id\"", 3,
		200,
	},
	{
		"search for cheese",
		"/search/cheese",
		"GET",
		``,
		"\\[\\]", 1,
		200,
	},
	{
		"get product 3",
		"/get/3",
		"GET",
		``,
		"orange", 1,
		200,
	},
	{
		"get on offer products",
		"/offers",
		"GET",
		``,
		"green", 1,
		200,
	},
	{
		"get all products",
		"/catalog",
		"GET",
		``,
		"\"id\"", 3,
		200,
	},
	{
		"get non-existent product",
		"/get/999",
		"GET",
		``,
		"", 0,
		404,
	},
}
