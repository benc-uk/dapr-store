package mock

import (
	"errors"
	"log"
	"time"

	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// OrderService mock
type OrderService struct {
}

var orders = map[string]spec.Order{
	"fake-order-01": {
		Title:   "A fake order",
		ID:      "fake-order-01",
		Items:   []string{"1", "3"},
		ForUser: "demo@example.net",
		Amount:  123.456,
		Status:  spec.OrderProcessing,
	},
	"fake-order-02": {
		Title:   "Another fake order",
		ID:      "fake-order-02",
		Items:   []string{"6", "2", "4"},
		ForUser: "test@example.net",
		Amount:  77.88,
		Status:  spec.OrderComplete,
	},
}

// GetOrder mock
func (s OrderService) GetOrder(orderID string) (*spec.Order, error) {
	order, exist := orders[orderID]
	if exist {
		return &order, nil
	}

	return nil, problem.New("err://not-found", "No data returned", 404, "Order: '"+orderID+"' not found", "orders")
}

// GetOrdersForUser mock
func (s OrderService) GetOrdersForUser(username string) ([]string, error) {
	return nil, nil
}

// ProcessOrder mock
func (s OrderService) ProcessOrder(order spec.Order) error {
	err := spec.Validate(order)
	if err != nil {
		return err
	}

	// Check we have a new order
	if order.Status != spec.OrderNew {
		return errors.New("Order not in correct status")
	}

	prob := s.AddOrder(order)
	if prob != nil {
		return prob
	}

	s.SetStatus(&order, spec.OrderReceived)

	log.Printf("### Order %s was saved to state store\n", order.ID)

	// Save order to blob storage as a text file "report"
	// Also email to the user via SendGrid
	// For these to work configure the components in cmd/orders/components
	// If un-configured then nothing happens, and no output is send or generated
	s.EmailNotify(order)
	s.SaveReport(order)

	// Fake background order processing
	time.AfterFunc(1*time.Second, func() {
		s.SetStatus(&order, spec.OrderProcessing)
	})

	// Fake background order completion
	time.AfterFunc(2*time.Second, func() {
		s.SetStatus(&order, spec.OrderComplete)
	})

	return nil
}

// AddOrder mock
func (s OrderService) AddOrder(order spec.Order) error {
	orders[order.ID] = order
	return nil
}

// SetStatus mock
func (s OrderService) SetStatus(order *spec.Order, status spec.OrderStatus) error {
	order.Status = status
	orders[order.ID] = *order
	return nil
}

// EmailNotify mock
func (s OrderService) EmailNotify(spec.Order) error {
	return nil
}

// SaveReport mock
func (s OrderService) SaveReport(spec.Order) error {
	return nil
}
