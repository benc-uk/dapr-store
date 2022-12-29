// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Main set of tests for users service and API
// ----------------------------------------------------------------------------

package main

import (
	"io"
	"log"
	"testing"

	"github.com/benc-uk/dapr-store/cmd/users/mock"
	"github.com/go-chi/chi/v5"

	"github.com/benc-uk/go-rest-api/pkg/api"
	"github.com/benc-uk/go-rest-api/pkg/auth"
	"github.com/benc-uk/go-rest-api/pkg/httptester"
)

func TestUsers(t *testing.T) {
	// Comment out to see logs
	log.SetOutput(io.Discard)

	// Mock of UserService
	mockUserSvc := &mock.UserService{}

	router := chi.NewRouter()
	api := API{
		api.NewBase("users", "ignore", "ignore", true),
		mockUserSvc,
	}

	api.addRoutes(router, auth.NewPassthroughValidator())

	httptester.Run(t, router, testCases)
}

var testCases = []httptester.TestCase{
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
		CheckBody:      "already",
		CheckBodyCount: 1,
		CheckStatus:    409,
	},
	{
		Name:   "register existing user 2",
		URL:    "/register",
		Method: "POST",
		Body: `{
			"username": "mock@example.net",
			"displayName": "Mock user"
		}`,
		CheckBody:      "already",
		CheckBodyCount: 1,
		CheckStatus:    409,
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
		CheckBody:      "invalid",
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
