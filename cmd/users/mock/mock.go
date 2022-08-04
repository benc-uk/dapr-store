// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Mock implementation of the UserService for testing
// ----------------------------------------------------------------------------

package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/benc-uk/dapr-store/cmd/users/spec"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// UserService mock
type UserService struct {
}

// Load mock data
var mockUsers []spec.User

func init() {
	mockJSON, err := ioutil.ReadFile("../../testing/mock-data/users.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(mockJSON, &mockUsers)
	if err != nil {
		panic(err)
	}
}

// GetUser mock
func (s UserService) GetUser(username string) (*spec.User, error) {
	for _, user := range mockUsers {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user %s not found", username)
}

// AddUser mock
func (s UserService) AddUser(user spec.User) error {
	userCheck, _ := s.GetUser(user.Username)
	if userCheck != nil {
		prob := problem.New("err://user-exists", user.Username+" already registered", 400, user.Username+" already registered", "users")
		return prob
	}

	mockUsers = append(mockUsers, user)

	return nil
}
