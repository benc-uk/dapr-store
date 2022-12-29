package mock

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/benc-uk/dapr-store/cmd/cart/impl"
	cartspec "github.com/benc-uk/dapr-store/cmd/cart/spec"
	orderspec "github.com/benc-uk/dapr-store/cmd/orders/spec"
)

// CartService mock
type CartService struct {
}

// Load mock data
var mockCarts []cartspec.Cart
var mockOrders []orderspec.Order

func init() {
	mockJSON, err := ioutil.ReadFile("../../testing/mock-data/carts.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(mockJSON, &mockCarts)
	if err != nil {
		panic(err)
	}

	mockJSON, err = ioutil.ReadFile("../../testing/mock-data/orders.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(mockJSON, &mockOrders)
	if err != nil {
		panic(err)
	}
}

// Get fetches saved cart for a given user, if not exists an empty cart is returned
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

// Submit a cart and turn into an order
func (s CartService) Submit(cart cartspec.Cart) (*orderspec.Order, error) {
	log.Printf("%+v", cart)

	if len(cart.Products) == 0 {
		return nil, impl.EmptyCartError()
	}

	return &mockOrders[0], nil
}

// SetProductCount updates the count of a given product in the cart
func (s CartService) SetProductCount(cart *cartspec.Cart, productID string, count int) error {
	if count < 0 {
		return impl.ProductCountError()
	}

	if count == 0 {
		delete(mockCarts[0].Products, productID)
		return nil
	}

	mockCarts[0].Products[productID] = count

	return nil
}

// Clear the cart
func (s CartService) Clear(cart *cartspec.Cart) error {
	cart.Products = map[string]int{}

	for i, c := range mockCarts {
		if c.ForUser == cart.ForUser {
			mockCarts[i] = *cart
		}
	}

	return nil
}
