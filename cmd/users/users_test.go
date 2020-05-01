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
