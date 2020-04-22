// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for users
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/benc-uk/dapr-store/pkg/auth"
	"github.com/benc-uk/dapr-store/pkg/models"
	"github.com/benc-uk/dapr-store/pkg/problem"
	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/register", auth.AuthMiddleware(api.registerUser)).Methods("POST")
	router.HandleFunc("/get/{username}", auth.AuthMiddleware(api.getUser))
	router.HandleFunc("/isregistered/{username}", api.checkRegistered)
}

//
// Register new user
//
func (api API) registerUser(resp http.ResponseWriter, req *http.Request) {
	cl, _ := strconv.Atoi(req.Header.Get("content-length"))
	if cl <= 0 {
		problem.New("err://body-missing", "Zero length body", 400, "Zero length body", daprHelper.AppInstanceName).Send(resp)
		return
	}

	user := models.User{}
	err := json.NewDecoder(req.Body).Decode(&user)

	// Some basic validation and checking on what we've been posted
	if err != nil {
		problem.New("err://json-decode", "Malformed user JSON", 400, "JSON could not be decoded", daprHelper.AppInstanceName).Send(resp)
		return
	}
	if len(user.DisplayName) == 0 || len(user.Username) == 0 {
		problem.New("err://json-error", "Malformed user JSON", 400, "User failed validation, check spec", daprHelper.AppInstanceName).Send(resp)
		return
	}
	log.Printf("### Registering user %+v\n", user)

	// Check is user already registered
	data, prob := daprHelper.GetState(daprStoreName, user.Username)
	if prob != nil {
		prob.Send(resp)
		return
	}
	log.Printf("### Existing user data %+v\n", string(data))
	if len(data) > 0 {
		problem.New("err://user-exists", user.Username+" already registered", 400, user.Username+" already registered", daprHelper.AppInstanceName).Send(resp)
		return
	}

	prob = daprHelper.SaveState(daprStoreName, user.Username, user)
	if prob != nil {
		prob.Send(resp)
		return
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
	data, prob := daprHelper.GetState(daprStoreName, vars["username"])
	if prob != nil {
		prob.Send(resp)
		return
	}

	if len(data) <= 0 {
		problem.New("err://not-found", "No data returned", 404, "Username: '"+vars["username"]+"' not found", daprHelper.AppInstanceName).Send(resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
}

//
// Returns 204 if registered and 404 if not
//
func (api API) checkRegistered(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	data, prob := daprHelper.GetState(daprStoreName, vars["username"])
	if prob != nil {
		prob.Send(resp)
		return
	}

	if len(data) <= 0 {
		resp.WriteHeader(404)
		return
	}

	resp.WriteHeader(204)
}
