// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Specification of the User entity and service
// ----------------------------------------------------------------------------

package spec

// A User holds information about a registered user
type User struct {
	UserID       string `json:"userId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	ProfileImage string `json:"profileImage"`
}

// UserService defines core CRUD methods a user service should have
type UserService interface {
	GetUser(userID string) (*User, error)
	AddUser(User) error
}
