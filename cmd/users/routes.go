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

	"github.com/benc-uk/dapr-store/cmd/users/spec"
	"github.com/benc-uk/dapr-store/pkg/auth"
	"github.com/benc-uk/dapr-store/pkg/problem"
	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/register", auth.JWTValidator(api.registerUser)).Methods("POST")
	router.HandleFunc("/get/{username}", auth.JWTValidator(api.getUser))
	router.HandleFunc("/isregistered/{username}", api.checkRegistered)
}

//
// Register new user
//
func (api API) registerUser(resp http.ResponseWriter, req *http.Request) {
	cl, _ := strconv.Atoi(req.Header.Get("content-length"))
	if cl <= 0 {
		problem.New("err://body-missing", "Zero length body", 400, "Zero length body", api.ServiceName).Send(resp)
		return
	}

	user := spec.User{}
	err := json.NewDecoder(req.Body).Decode(&user)

	// Some basic validation and checking on what we've been posted
	if err != nil {
		problem.New("err://json-decode", "Malformed user JSON", 400, "JSON could not be decoded", api.ServiceName).Send(resp)
		return
	}

	if len(user.DisplayName) == 0 || len(user.Username) == 0 {
		problem.New("err://json-error", "Malformed user JSON", 400, "User failed validation, check spec", api.ServiceName).Send(resp)
		return
	}

	log.Printf("### Registering user %+v\n", user)

	err = api.service.AddUser(user)
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)

		return
	}

	// Send success message back
	resp.Header().Set("Content-Type", "application/json")
	_, _ = resp.Write([]byte(fmt.Sprintf(`{"registrationStatus":"success", "username":"%s"}`, user.Username)))
}

//
// Fetch existing user, return 404 if they don't exist
//
func (api API) getUser(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	user, err := api.service.GetUser(vars["username"])
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)

		return
	}

	if user == nil {
		resp.WriteHeader(404)
		return
	}

	resp.Header().Set("Content-Type", "application/json")

	json, _ := json.Marshal(user)
	_, _ = resp.Write(json)
}

//
// Returns 204 if registered and 404 if not
//
func (api API) checkRegistered(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	_, err := api.service.GetUser(vars["username"])
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)

		return
	}

	resp.WriteHeader(204)
}
