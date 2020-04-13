// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for orders
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/benc-uk/dapr-store/common"
	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/register", api.registerUser).Methods("POST")
	router.HandleFunc("/get/{username}", api.getUser)
}

//
// Register new user
//
func (api API) registerUser(resp http.ResponseWriter, req *http.Request) {
	cl, _ := strconv.Atoi(req.Header.Get("content-length"))
	if cl <= 0 {
		common.Problem{"json-error", "Zero length body", 400, "Body is required", serviceName}.HttpSend(resp)
		return
	}

	user := common.User{}
	err := json.NewDecoder(req.Body).Decode(&user)

	// Some basic validation and checking on what we've been posted
	if err != nil {
		common.Problem{"json-error", "JSON decoding error", 400, err.Error(), serviceName}.HttpSend(resp)
		return
	}
	if len(user.DisplayName) == 0 || len(user.Username) == 0 {
		common.Problem{"json-error", "Malformed user JSON", 400, "Validation failed, check user schema", serviceName}.HttpSend(resp)
		return
	}

	// Check is user already registered
	data, err := common.GetState(resp, daprPort, daprStateStore, serviceName, user.Username)
	if len(data) > 0 {
		common.Problem{"user-exists", user.Username + " already registered", 400, "User is already registered!", serviceName}.HttpSend(resp)
		return
	}

	err = common.SaveState(resp, daprPort, daprStateStore, serviceName, user.Username, user)
	if err != nil {
		return // Error will have already been written to resp
	}

	// Send success message back
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(fmt.Sprintf(`{"registrationStatus":"success", "username":"%s"`, user.Username)))
}

//
// Fetch existing user
//
func (api API) getUser(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	data, err := common.GetState(resp, daprPort, daprStateStore, serviceName, vars["username"])

	if len(data) <= 0 {
		common.Problem{"user-not-found", vars["username"] + " not found", 404, "User is not registered", serviceName}.HttpSend(resp)
		return
	}

	if err != nil {
		return // Error will have already been written to resp
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
}
