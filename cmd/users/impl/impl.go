// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr based implementation of the UserService
// ----------------------------------------------------------------------------

package impl

import (
	"context"
	"encoding/json"
	"log"

	"github.com/benc-uk/dapr-store/cmd/users/spec"

	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/benc-uk/dapr-store/pkg/problem"
	dapr "github.com/dapr/go-sdk/client"
)

// UserService is a Dapr based implementation of UserService interface
type UserService struct {
	storeName   string
	serviceName string
	client      dapr.Client
}

// NewService creates a new UserService
func NewService(serviceName string) *UserService {
	storeName := env.GetEnvString("DAPR_STORE_NAME", "statestore")

	// Set up Dapr client & checks for Dapr sidecar, otherwise die
	client, err := dapr.NewClient()
	if err != nil {
		log.Panicln("FATAL! Dapr process/sidecar NOT found. Terminating!")
	}

	return &UserService{
		storeName,
		serviceName,
		client,
	}
}

// AddUser registers a new user and stores in Dapr state
func (s *UserService) AddUser(user spec.User) error {
	// Check is user already registered
	data, err := s.client.GetState(context.Background(), s.storeName, user.Username, nil)
	if err != nil {
		return problem.NewDaprStateProblem(err, s.serviceName)
	}

	// If we get any data, that means we found a user, that's an error in our case
	if data.Value != nil {
		prob := problem.New("err://user-exists", user.Username+" already registered", 400, user.Username+" already registered", s.serviceName)
		return prob
	}

	// Call Dapr client to save state
	jsonPayload, err := json.Marshal(user)
	if err != nil {
		return problem.New500("err://json-marshall", "State JSON marshalling error", s.serviceName, nil, err)
	}
	if err := s.client.SaveState(context.Background(), s.storeName, user.Username, jsonPayload, nil); err != nil {
		return problem.NewDaprStateProblem(err, s.serviceName)
	}

	return nil
}

// GetUser fetches a user from Dapr state
func (s *UserService) GetUser(username string) (*spec.User, error) {
	data, err := s.client.GetState(context.Background(), s.serviceName, username, nil)
	if err != nil {
		return nil, problem.NewDaprStateProblem(err, s.serviceName)
	}

	if data.Value == nil {
		prob := problem.New("err://not-found", "No data returned", 404, "Username: '"+username+"' not found", s.serviceName)
		return nil, prob
	}

	user := &spec.User{}
	err = json.Unmarshal(data.Value, user)
	if err != nil {
		prob := problem.New("err://json-decode", "Malformed user JSON", 500, "JSON could not be decoded", s.serviceName)
		return nil, prob
	}

	return user, nil
}
