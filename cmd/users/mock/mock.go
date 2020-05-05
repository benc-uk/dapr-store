// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Mock implementation of the UserService for testing
// ----------------------------------------------------------------------------

package mock

import (
	"encoding/json"
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
	mockJson, err := ioutil.ReadFile("../../etc/mock-data/users.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(mockJson, &mockUsers)
}

// GetUser mock
func (s UserService) GetUser(username string) (*spec.User, error) {
	for _, user := range mockUsers {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, nil
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
