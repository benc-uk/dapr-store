package spec

import "github.com/benc-uk/dapr-store/cmd/orders/spec"

// OrderService defines core CRUD methods a orders service should have
type CartService interface {
	SubmitOrder(*spec.Order) error
}
