package impl

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	cartspec "github.com/benc-uk/dapr-store/cmd/cart/spec"
	orderspec "github.com/benc-uk/dapr-store/cmd/orders/spec"
	productspec "github.com/benc-uk/dapr-store/cmd/products/spec"
	"github.com/benc-uk/dapr-store/pkg/dapr"
	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// CartService is a Dapr implementation of CartService interface
type CartService struct {
	*dapr.Helper
	topicName string
	storeName string
}

//
// NewService creates a new OrderService
//
func NewService(serviceName string) *CartService {
	// Set up Dapr & checks for Dapr sidecar port, abort
	helper := dapr.NewHelper(serviceName)
	if helper == nil {
		os.Exit(1)
	}
	topicName := env.GetEnvString("DAPR_ORDERS_TOPIC", "orders-queue")
	storeName := env.GetEnvString("DAPR_STORE_NAME", "statestore")

	log.Printf("### Dapr pub/sub topic name: %s\n", topicName)
	log.Printf("### Dapr state store name:   %s\n", storeName)

	return &CartService{
		helper,
		topicName,
		storeName,
	}
}

//
// Get fetches saved cart for a given user, if not exists an empty cart is returned
//
func (s CartService) Get(username string) (*cartspec.Cart, error) {
	data, prob := s.GetState(s.storeName, username)
	if prob != nil {
		return nil, prob
	}

	if len(data) <= 0 {
		cart := &cartspec.Cart{}
		cart.ForUser = username
		cart.Products = make(map[string]int)
		return cart, nil
	}

	cart := &cartspec.Cart{}
	err := json.Unmarshal(data, cart)
	if err != nil {
		prob := problem.New("err://json-decode", "Malformed cart JSON", 500, "JSON could not be decoded", s.ServiceName)
		return nil, prob
	}

	return cart, nil
}

//
// Submit a cart and turn into an order
//
func (s CartService) Submit(cart cartspec.Cart) (*orderspec.Order, error) {
	if len(cart.Products) == 0 {
		return nil, problem.New("err://order-cart", "Cart empty", 400, "No items in cart", s.ServiceName)
	}

	// Build up line item array
	lineItems := []orderspec.LineItem{}

	var orderAmount float32 = 0.0
	for productID, count := range cart.Products {
		resp, err := s.InvokeGet("products", `get/`+productID)
		if err != nil || resp.StatusCode != 200 {
			return nil, problem.NewAPIProblem("err://cart-product", "Submit cart, product lookup error "+productID, s.ServiceName, resp, err)
		}

		product := &productspec.Product{}
		err = json.NewDecoder(resp.Body).Decode(product)
		if err != nil {
			prob := problem.New("err://json-decode", "Malformed JSON", 500, "Product JSON could not be decoded", s.ServiceName)
			return nil, prob
		}
		lineItem := &orderspec.LineItem{
			Product: *product,
			Count:   count,
		}
		lineItems = append(lineItems, *lineItem)

		orderAmount += (product.Cost * float32(count))
	}

	order := &orderspec.Order{
		Title:     "Order " + time.Now().Format("15:04 Jan 2 2006"),
		Amount:    orderAmount,
		ForUser:   cart.ForUser,
		ID:        makeID(5),
		Status:    orderspec.OrderNew,
		LineItems: lineItems,
	}

	prob := s.PublishMessage(s.topicName, order)
	if prob != nil {
		return nil, prob
	}
	s.Clear(&cart)

	return order, nil
}

//
// SetProductCount updates the count of a given product in the cart
//
func (s CartService) SetProductCount(cart *cartspec.Cart, productID string, count int) error {
	if count < 0 {
		return problem.New("err://invalid-request", "SetProductCount error", 400, "Count can not be negative", s.ServiceName)
	}

	if count == 0 {
		delete(cart.Products, productID)
	} else {
		cart.Products[productID] = count
	}

	prob := s.SaveState(s.storeName, cart.ForUser, cart)
	if prob != nil {
		return prob
	}

	return nil
}

//
// Clear the cart
//
func (s CartService) Clear(cart *cartspec.Cart) error {
	cart.Products = map[string]int{}
	prob := s.SaveState(s.storeName, cart.ForUser, cart)
	if prob != nil {
		return prob
	}
	return nil
}

//
// Scummy but functional ID generator
//
func makeID(length int) string {
	id := ""
	possible := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		id += string(possible[rand.Intn(len(possible)-1)])
	}

	return id
}
