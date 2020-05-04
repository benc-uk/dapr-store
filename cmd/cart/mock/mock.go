package mock

import (
	cart "github.com/benc-uk/dapr-store/cmd/cart/spec"
	orderspec "github.com/benc-uk/dapr-store/cmd/orders/spec"
	productspec "github.com/benc-uk/dapr-store/cmd/products/spec"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// CartService mock
type CartService struct {
}

var mockCart = &cart.Cart{
	Products: map[string]int{},
	ForUser:  "demo@example.net",
}

func (s CartService) Get(username string) (*cart.Cart, error) {
	if username != "demo@example.net" {
		cart := &cart.Cart{}
		cart.ForUser = username
		cart.Products = make(map[string]int)
		return cart, nil
	}
	return mockCart, nil
}

func (s CartService) Submit(cart cart.Cart) (*orderspec.Order, error) {
	if len(cart.Products) == 0 {
		return nil, problem.New("err://bad", "Cart empty", 400, "Cart empty", "mock-cart")
	}

	o := &orderspec.Order{
		Title:   "Mock Order",
		Amount:  12.34,
		ForUser: cart.ForUser,
		ID:      "order-01",
		Status:  orderspec.OrderNew,
		LineItems: []orderspec.LineItem{
			{
				Count: 1,
				Product: productspec.Product{
					ID:          "4",
					Name:        "foo",
					Cost:        12.34,
					Description: "blah",
					Image:       "blah.jpg",
					OnOffer:     false,
				},
			},
		},
	}
	return o, nil
}

func (s CartService) SetProductCount(cart *cart.Cart, productId string, count int) error {
	if count < 0 {
		return problem.New("err://bad", "SetProductCount", 500, "count can not be negative", "mock-cart")
	}
	if count == 0 {
		delete(mockCart.Products, productId)
		return nil
	}
	mockCart.Products[productId] = count
	return nil
}

func (s CartService) Clear(cart *cart.Cart) error {
	cart.Products = map[string]int{}
	return nil
}
