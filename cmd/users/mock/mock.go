// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Mock implementation of the UserService for testing
// ----------------------------------------------------------------------------

package mock

import (
	"encoding/json"
	"os"

	"github.com/benc-uk/dapr-store/cmd/users/impl"
	"github.com/benc-uk/dapr-store/cmd/users/spec"
)

// UserService mock
type UserService struct {
}

// Load mock data
var mockUsers []spec.User

func init() {
	mockJSON, err := os.ReadFile("../../testing/mock-data/users.json")
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

	return nil, impl.UserNotFoundError()
}

// AddUser mock
func (s UserService) AddUser(user spec.User) error {
	userCheck, _ := s.GetUser(user.Username)
	if userCheck != nil {
		return impl.UserDuplicateError()
	}

	mockUsers = append(mockUsers, user)

	return nil
}
