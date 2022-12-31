// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr based implementation of the CartService
// ----------------------------------------------------------------------------

package impl

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	cartspec "github.com/benc-uk/dapr-store/cmd/cart/spec"
	orderspec "github.com/benc-uk/dapr-store/cmd/orders/spec"
	productspec "github.com/benc-uk/dapr-store/cmd/products/spec"

	"github.com/benc-uk/go-rest-api/pkg/env"
	dapr "github.com/dapr/go-sdk/client"
)

// CartService is a Dapr implementation of CartService interface
type CartService struct {
	pubSubName  string // Name of Papr pub/sub component for orders
	topicName   string // Name of Dapr pub/sub topic for orders
	storeName   string // Name of Dapr state store
	serviceName string
	client      dapr.Client
}

// NewService creates a new CartService
func NewService(serviceName string) *CartService {
	topicName := env.GetEnvString("DAPR_ORDERS_TOPIC", "orders-queue")
	storeName := env.GetEnvString("DAPR_STORE_NAME", "statestore")
	pubSubName := env.GetEnvString("DAPR_PUBSUB_NAME", "pubsub")

	// Set up Dapr client & checks for Dapr sidecar, otherwise die
	client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("FATAL! Dapr process/sidecar NOT found. Terminating!")
	}

	return &CartService{
		pubSubName,
		topicName,
		storeName,
		serviceName,
		client,
	}
}

// Get fetches saved cart for a given user, if none exists an empty cart is always returned
func (s CartService) Get(username string) (*cartspec.Cart, error) {
	data, err := s.client.GetState(context.Background(), s.storeName, username, nil)
	if err != nil {
		return nil, err
	}

	// Create an empty cart
	if data.Value == nil {
		cart := &cartspec.Cart{}
		cart.ForUser = username
		cart.Products = make(map[string]int)

		return cart, nil
	}

	cart := &cartspec.Cart{}

	if err = json.Unmarshal(data.Value, cart); err != nil {
		// The cart is somehow corrupt - remove it from state, or we'll get stuck
		_ = s.client.DeleteState(context.Background(), s.storeName, username, nil)

		log.Printf("### Warning: Corrupt cart for user %s was removed!!", username)

		cart.ForUser = username
		cart.Products = make(map[string]int)
	}

	return cart, nil
}

// Submit a cart and turn into an order
func (s CartService) Submit(cart cartspec.Cart) (*orderspec.Order, error) {
	if len(cart.Products) == 0 {
		return nil, EmptyCartError()
	}

	// Build up line item array
	lineItems := []orderspec.LineItem{}

	// Process the cart server side, calculating the order price
	// This involves a service to service call to invoke the products service
	var orderAmount float32

	for productID, count := range cart.Products {
		resp, err := s.client.InvokeMethod(context.Background(), "products", `get/`+productID, "get")
		if err != nil {
			return nil, ProductLookupError(productID)
		}

		product := &productspec.Product{}

		err = json.Unmarshal(resp, product)
		if err != nil {
			return nil, err
		}

		lineItem := &orderspec.LineItem{
			Product: *product,
			Count:   count,
		}
		lineItems = append(lineItems, *lineItem)

		orderAmount += (product.Cost * float32(count))
	}

	// Publish order to the orders queue
	order := &orderspec.Order{
		Title:     "Order " + time.Now().Format("15:04 Jan 2 2006"),
		Amount:    orderAmount,
		ForUser:   cart.ForUser,
		ID:        makeID(5),
		Status:    orderspec.OrderNew,
		LineItems: lineItems,
	}

	err := s.client.PublishEvent(context.Background(), s.pubSubName, s.topicName, order)
	if err != nil {
		return nil, err
	}

	err = s.Clear(&cart)
	if err != nil {
		// Log but don't return the error, as the order was published
		log.Printf("### Warning failed to clear cart %s", err)
	}

	return order, nil
}

// SetProductCount updates the count of a given product in the cart
func (s CartService) SetProductCount(cart *cartspec.Cart, productID string, count int) error {
	if count < 0 {
		return ProductCountError()
	}

	if count == 0 {
		delete(cart.Products, productID)
	} else {
		cart.Products[productID] = count
	}

	// Call Dapr client to save state
	jsonPayload, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	if err = s.client.SaveState(context.Background(), s.storeName, cart.ForUser, jsonPayload, nil); err != nil {
		return err
	}

	return nil
}

// Clear the cart
func (s CartService) Clear(cart *cartspec.Cart) error {
	cart.Products = map[string]int{}
	// Call Dapr client to save state
	jsonPayload, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	if err = s.client.SaveState(context.Background(), s.storeName, cart.ForUser, jsonPayload, nil); err != nil {
		return err
	}

	return nil
}

// Scummy but functional ID generator
func makeID(length int) string {
	id := ""
	possible := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		id += string(possible[rand.Intn(len(possible)-1)])
	}

	return id
}
