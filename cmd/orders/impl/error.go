// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Errors supporting the orders service
// ----------------------------------------------------------------------------

package impl

const NotFoundError = "order not found"
const StatusError = "order status invalid"

type OrdersError struct {
	err string
}

func (e OrdersError) Error() string {
	return e.err
}

func OrderNotFoundError() OrdersError {
	return OrdersError{NotFoundError}
}

func OrderStatusError() OrdersError {
	return OrdersError{NotFoundError}
}
