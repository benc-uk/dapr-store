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
	"github.com/benc-uk/go-rest-api/pkg/env"

	dapr "github.com/dapr/go-sdk/client"
)

// UserService is a Dapr based implementation of UserService interface
type UserService struct {
	storeName   string // Name of Dapr state store
	serviceName string
	client      dapr.Client
}

// NewService creates a new UserService
func NewService(serviceName string) *UserService {
	storeName := env.GetEnvString("DAPR_STORE_NAME", "statestore")

	// Set up Dapr client & checks for Dapr sidecar, otherwise die
	client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("FATAL! Dapr process/sidecar NOT found. Terminating!")
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
	data, err := s.client.GetState(context.Background(), s.storeName, user.UserID, nil)
	if err != nil {
		return err
	}

	// If we get any data, that means we found a user, that's an error in our case
	if data.Value != nil {
		return UserDuplicateError()
	}

	// Call Dapr client to save state
	jsonPayload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := s.client.SaveState(context.Background(), s.storeName, user.UserID, jsonPayload, nil); err != nil {
		return err
	}

	return nil
}

// GetUser fetches a user from Dapr state
func (s *UserService) GetUser(userID string) (*spec.User, error) {
	data, err := s.client.GetState(context.Background(), s.storeName, userID, nil)
	if err != nil {
		return nil, err
	}

	if data.Value == nil {
		return nil, UserNotFoundError()
	}

	user := &spec.User{}

	if err = json.Unmarshal(data.Value, user); err != nil {
		return nil, err
	}

	return user, nil
}
