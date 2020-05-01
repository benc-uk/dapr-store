// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Mock implementation of the UserService for testing
// ----------------------------------------------------------------------------

package mock

import (
	"github.com/benc-uk/dapr-store/cmd/users/spec"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// UserService mock
type UserService struct {
}

var users = []spec.User{
	{
		DisplayName:  "Mock user",
		Username:     "demo@example.net",
		ProfileImage: "face.jpg",
	},
}

// GetUser mock
func (s UserService) GetUser(username string) (*spec.User, error) {
	for _, user := range users {
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
	users = append(users, user)
	return nil
}
