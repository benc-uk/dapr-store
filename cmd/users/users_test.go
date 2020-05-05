// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Main set of tests for users service and API
// ----------------------------------------------------------------------------

package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/benc-uk/dapr-store/cmd/users/mock"
	"github.com/benc-uk/dapr-store/pkg/api"
	"github.com/benc-uk/dapr-store/pkg/apitests"
	"github.com/gorilla/mux"
)

func TestUsers(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	// Mock of UserService
	mockUserSvc := &mock.UserService{}

	router := mux.NewRouter()
	api := API{
		api.NewBase("users", "ignore", "ignore", true, router),
		mockUserSvc,
	}
	api.addRoutes(router)

	apitests.Run(t, router, testCases)
}

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
			"username": "mock@example.net",
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
