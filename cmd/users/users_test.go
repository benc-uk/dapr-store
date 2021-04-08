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
		Name:   "register valid user",
		URL:    "/register",
		Method: "POST",
		Body: `{
			"username": "test@example.net",
			"displayName": "Mr Test"
		}`,
		CheckBody:      "",
		CheckBodyCount: 0,
		CheckStatus:    200,
	},
	{
		Name:   "register existing user 1",
		URL:    "/register",
		Method: "POST",
		Body: `{
			"username": "test@example.net",
			"displayName": "Mr Test"
		}`,
		CheckBody:      "already registered",
		CheckBodyCount: 2,
		CheckStatus:    400,
	},
	{
		Name:   "register existing user 2",
		URL:    "/register",
		Method: "POST",
		Body: `{
			"username": "mock@example.net",
			"displayName": "Mock user"
		}`,
		CheckBody:      "already registered",
		CheckBodyCount: 2,
		CheckStatus:    400,
	},
	{
		Name:   "register invalid user",
		URL:    "/register",
		Method: "POST",
		Body: `{
			"displayName": "Mr Test"
		}`,
		CheckBody:      "failed validation",
		CheckBodyCount: 1,
		CheckStatus:    400,
	},
	{
		Name:           "register invalid json",
		URL:            "/register",
		Method:         "POST",
		Body:           `lemon_curd`,
		CheckBody:      "Malformed",
		CheckBodyCount: 1,
		CheckStatus:    400,
	},
	{
		Name:           "get existing user",
		URL:            "/get/test@example.net",
		Method:         "GET",
		Body:           ``,
		CheckBody:      "Mr Test",
		CheckBodyCount: 1,
		CheckStatus:    200,
	},
	{
		Name:           "get non-registered user",
		URL:            "/get/idontexist@example.net",
		Method:         "GET",
		Body:           ``,
		CheckBody:      "",
		CheckBodyCount: 0,
		CheckStatus:    404,
	},
}
