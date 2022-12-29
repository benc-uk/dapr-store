// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Errors supporting the UserService
// ----------------------------------------------------------------------------

package impl

const NotFoundError = "user not found"
const DuplicateError = "user already exists"

type UserError struct {
	err string
}

func (e UserError) Error() string {
	return e.err
}

func UserNotFoundError() UserError {
	return UserError{NotFoundError}
}

func UserDuplicateError() UserError {
	return UserError{DuplicateError}
}
