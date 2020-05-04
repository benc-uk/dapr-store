package mock

import (
	"errors"
	"log"
	"time"

	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	products "github.com/benc-uk/dapr-store/cmd/products/spec"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// OrderService mock
type OrderService struct {
}

var Orders = map[string]spec.Order{
	"fake-order-01": {
		Title: "A fake order",
		ID:    "fake-order-01",
		LineItems: []spec.LineItem{
			{
				Count: 1,
				Product: products.Product{
					ID:          "4",
					Name:        "foo",
					Cost:        12.34,
					Description: "blah",
					Image:       "fo.jpg",
					OnOffer:     false,
				},
			},
		},
		ForUser: "demo@example.net",
		Amount:  123.456,
		Status:  spec.OrderNew,
	},

	"fake-order-02": {
		Title: "Another fake order",
		ID:    "fake-order-02",
		LineItems: []spec.LineItem{
			{
				Count: 2,
				Product: products.Product{
					ID:          "7",
					Name:        "bar",
					Cost:        88.30,
					Description: "bar blah",
					Image:       "bar.jpg",
					OnOffer:     true,
				},
			},
		},
		ForUser: "test@example.net",
		Amount:  77.88,
		Status:  spec.OrderComplete,
	},
}

// GetOrder mock
func (s OrderService) GetOrder(orderID string) (*spec.Order, error) {
	order, exist := Orders[orderID]
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
	Orders[order.ID] = order
	return nil
}

// SetStatus mock
func (s OrderService) SetStatus(order *spec.Order, status spec.OrderStatus) error {
	order.Status = status
	Orders[order.ID] = *order
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
