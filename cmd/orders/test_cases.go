// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Test cases for testing orders API, run by orders_test.go
// ----------------------------------------------------------------------------

package main

import "github.com/benc-uk/dapr-store/pkg/apitests"

var testCases = []apitests.Test{
	{
		"get an existing order",
		"/get/fake-order-01",
		"GET",
		``,
		"fake", 1,
		200,
	},
	{
		"get an non-existent order",
		"/get/foo",
		"GET",
		``,
		"not found", 1,
		404,
	},
}
