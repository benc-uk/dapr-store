package mock

import (
	"encoding/json"
	"io/ioutil"
	"log"

	cartspec "github.com/benc-uk/dapr-store/cmd/cart/spec"
	orderspec "github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// CartService mock
type CartService struct {
}

// Load mock data
var mockCarts []cartspec.Cart
var mockOrders []orderspec.Order

func init() {
	mockJSON, err := ioutil.ReadFile("../../etc/mock-data/carts.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(mockJSON, &mockCarts)
	mockJSON, err = ioutil.ReadFile("../../etc/mock-data/orders.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(mockJSON, &mockOrders)
}

//
// Get fetches saved cart for a given user, if not exists an empty cart is returned
//
func (s CartService) Get(username string) (*cartspec.Cart, error) {
	for _, cart := range mockCarts {
		if cart.ForUser == username {
			return &cart, nil
		}
	}

	cart := &cartspec.Cart{}
	cart.ForUser = username
	cart.Products = make(map[string]int)
	return cart, nil
}

//
// Submit a cart and turn into an order
//
func (s CartService) Submit(cart cartspec.Cart) (*orderspec.Order, error) {
	log.Printf("%+v", cart)
	if len(cart.Products) == 0 {
		return nil, problem.New("err://bad", "Cart empty", 400, "Cart empty", "mock-cart")
	}

	return &mockOrders[0], nil
}

//
// SetProductCount updates the count of a given product in the cart
//
func (s CartService) SetProductCount(cart *cartspec.Cart, productID string, count int) error {
	if count < 0 {
		return problem.New("err://bad", "SetProductCount", 500, "count can not be negative", "mock-cart")
	}
	if count == 0 {
		delete(mockCarts[0].Products, productID)
		return nil
	}
	mockCarts[0].Products[productID] = count
	return nil
}

//
// Clear the cart
//
func (s CartService) Clear(cart *cartspec.Cart) error {
	cart.Products = map[string]int{}
	for i, c := range mockCarts {
		if c.ForUser == cart.ForUser {
			mockCarts[i] = *cart
		}
	}
	return nil
}
