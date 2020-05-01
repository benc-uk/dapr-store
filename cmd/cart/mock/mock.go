package mock

import "github.com/benc-uk/dapr-store/cmd/orders/spec"

// CartService mock
type CartService struct {
}

// SubmitOrder does mock order submission
func (s CartService) SubmitOrder(order *spec.Order) error {
	order.ID = "fake_id_01"
	return nil
}
