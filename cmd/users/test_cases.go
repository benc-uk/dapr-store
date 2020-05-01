// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Test cases for testing user API, run by users_test.go
// ----------------------------------------------------------------------------

package main

import "github.com/benc-uk/dapr-store/pkg/apitests"

var testCases = []apitests.Test{
	{
		"register valid user",
		"/register",
		"POST",
		`{
			"username": "test@example.net",
			"displayName": "Mr Test"
		}`,
		"", 0,
		200,
	},
	{
		"register existing user 1",
		"/register",
		"POST",
		`{
			"username": "test@example.net",
			"displayName": "Mr Test"
		}`,
		"already registered", 2,
		400,
	},
	{
		"register existing user 2",
		"/register",
		"POST",
		`{
			"username": "demo@example.net",
			"displayName": "Mock user"
		}`,
		"already registered", 2,
		400,
	},
	{
		"register invalid user",
		"/register",
		"POST",
		`{
			"displayName": "Mr Test"
		}`,
		"failed validation", 1,
		400,
	},
	{
		"register invalid json",
		"/register",
		"POST",
		`lemon_curd`,
		"Malformed", 1,
		400,
	},
	{
		"get existing user",
		"/get/test@example.net",
		"GET",
		``,
		"Mr Test", 1,
		200,
	},
	{
		"get non-registered user",
		"/get/idontexist@example.net",
		"GET",
		``,
		"", 0,
		404,
	},
}
