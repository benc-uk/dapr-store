// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for users
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/benc-uk/dapr-store/cmd/users/impl"
	"github.com/benc-uk/dapr-store/cmd/users/spec"
	"github.com/benc-uk/go-rest-api/pkg/auth"
	"github.com/benc-uk/go-rest-api/pkg/problem"
	"github.com/go-chi/chi/v5"
)

// All routes we need should be registered here
func (api API) addRoutes(router chi.Router, v auth.Validator) {
	router.Post("/register", v.Protect(api.registerUser))
	router.Get("/get/{username}", v.Protect(api.getUser))
	router.Get("/isregistered/{username}", api.checkRegistered)
}

// Register new user
func (api API) registerUser(resp http.ResponseWriter, req *http.Request) {
	user := spec.User{}

	// Some basic validation and checking on what we've been posted
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		problem.Wrap(400, req.RequestURI, "new-user", err).Send(resp)
		return
	}

	if len(user.DisplayName) == 0 || len(user.Username) == 0 {
		problem.Wrap(400, req.RequestURI, "new-user", errors.New("failed validation, check spec")).Send(resp)
		return
	}

	log.Printf("### Registering user %+v\n", user)

	if err := api.service.AddUser(user); err != nil {
		if userError, isError := err.(impl.UserError); isError && userError.Error() == impl.DuplicateError {
			problem.Wrap(409, req.RequestURI, user.Username, err).Send(resp)
			return
		}

		problem.Wrap(500, req.RequestURI, user.Username, err).Send(resp)

		return
	}

	// Send success message back
	api.ReturnJSON(resp, map[string]string{
		"registrationStatus": "success",
		"username":           user.Username,
	})
}

// Fetch existing user, return 404 if they don't exist
func (api API) getUser(resp http.ResponseWriter, req *http.Request) {
	username := chi.URLParam(req, "username")

	user, err := api.service.GetUser(username)
	if err != nil {
		if userError, isError := err.(impl.UserError); isError && userError.Error() == impl.NotFoundError {
			problem.Wrap(404, req.RequestURI, username, err).Send(resp)
			return
		}

		problem.Wrap(500, req.RequestURI, username, err).Send(resp)

		return
	}

	api.ReturnJSON(resp, user)
}

// Returns 204 if registered and 404 if not
func (api API) checkRegistered(resp http.ResponseWriter, req *http.Request) {
	username := chi.URLParam(req, "username")

	if _, err := api.service.GetUser(username); err != nil {
		if userError, isError := err.(impl.UserError); isError && userError.Error() == impl.NotFoundError {
			problem.Wrap(404, req.RequestURI, username, err).Send(resp)
			return
		}

		problem.Wrap(500, req.RequestURI, username, err).Send(resp)

		return
	}

	resp.WriteHeader(204)
}
