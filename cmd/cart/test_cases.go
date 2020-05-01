// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Test cases for testing cart API, run by cart_test.go
// ----------------------------------------------------------------------------

package main

import "github.com/benc-uk/dapr-store/pkg/apitests"

var testCases = []apitests.Test{
	{
		"submit valid order",
		"/submit",
		"POST",
		`{
			"title": "test order", 
			"amount": 11.22,
			"forUser": "demo@example.net",
			"items": ["1", "7"]
		}`,
		"fake_id_01", 1,
		200,
	},
	{
		"submit invalid order",
		"/submit",
		"POST",
		`{
			"title": "test order", 
			"amount": 666
		}`,
		"Order failed validation", 1,
		400,
	},
	{
		"submit invalid json",
		"/submit",
		"POST",
		`blah_blah`,
		"Malformed", 1,
		400,
	},
	{
		"submit empty order",
		"/submit",
		"POST",
		``,
		"Zero length body", 1,
		400,
	},
	{
		"submit POST only",
		"/submit",
		"GET",
		``,
		"", 0,
		405,
	},
}
