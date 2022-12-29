// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Errors supporting the cart service
// ----------------------------------------------------------------------------

package impl

const EmptyError = "cart is empty"
const CountError = "product count must be > 0"
const LookupError = "product lookup failed: "

type CartError struct {
	err string
}

func (e CartError) Error() string {
	return e.err
}

func EmptyCartError() CartError {
	return CartError{EmptyError}
}

func ProductCountError() CartError {
	return CartError{CountError}
}

func ProductLookupError(prodID string) CartError {
	return CartError{LookupError + prodID}
}
