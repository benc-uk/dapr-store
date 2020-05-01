package spec

import "github.com/benc-uk/dapr-store/cmd/orders/spec"

// Cart isn't currently used
type Cart struct {
	items []string
}

// CartService defines core CRUD methods a cart service should have
type CartService interface {
	SubmitOrder(*spec.Order) error
}
