package mock

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"time"

	orderspec "github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// OrderService mock
type OrderService struct {
}

// Load mock data
var MockOrders []orderspec.Order
var mockUserOrders []string

func init() {
	mockJson, err := ioutil.ReadFile("../../etc/mock-data/orders.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(mockJson, &MockOrders)
	mockJson, err = ioutil.ReadFile("../../etc/mock-data/user-orders.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(mockJson, &mockUserOrders)
}

// GetOrder mock
func (s OrderService) GetOrder(orderID string) (*orderspec.Order, error) {
	if orderID == MockOrders[0].ID {
		return &MockOrders[0], nil
	}

	return nil, problem.New("err://not-found", "No data returned", 404, "Order: '"+orderID+"' not found", "orders")
}

// GetOrdersForUser mock
func (s OrderService) GetOrdersForUser(username string) ([]string, error) {
	return nil, nil
}

// ProcessOrder mock
func (s OrderService) ProcessOrder(order orderspec.Order) error {
	err := orderspec.Validate(order)
	if err != nil {
		return err
	}

	// Check we have a new order
	if order.Status != orderspec.OrderNew {
		return errors.New("Order not in correct status")
	}

	prob := s.AddOrder(order)
	if prob != nil {
		return prob
	}

	s.SetStatus(&order, orderspec.OrderReceived)

	log.Printf("### Order %s was saved to state store\n", order.ID)

	// Save order to blob storage as a text file "report"
	// Also email to the user via SendGrid
	// For these to work configure the components in cmd/orders/components
	// If un-configured then nothing happens, and no output is send or generated
	s.EmailNotify(order)
	s.SaveReport(order)

	// Fake background order processing
	time.AfterFunc(1*time.Second, func() {
		s.SetStatus(&order, orderspec.OrderProcessing)
	})

	// Fake background order completion
	time.AfterFunc(2*time.Second, func() {
		s.SetStatus(&order, orderspec.OrderComplete)
	})

	return nil
}

// AddOrder mock
func (s OrderService) AddOrder(order orderspec.Order) error {
	MockOrders = append(MockOrders, order)
	return nil
}

// SetStatus mock
func (s OrderService) SetStatus(order *orderspec.Order, status orderspec.OrderStatus) error {
	order.Status = status
	MockOrders[0] = *order
	return nil
}

// EmailNotify mock
func (s OrderService) EmailNotify(orderspec.Order) error {
	return nil
}

// SaveReport mock
func (s OrderService) SaveReport(orderspec.Order) error {
	return nil
}
