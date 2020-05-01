package spec

import (
	"errors"
)

// Order holds information about a customer order
type Order struct {
	ID      string      `json:"id"`
	Title   string      `json:"title"`
	Amount  float32     `json:"amount"`
	Items   []string    `json:"items"` // List of Product.ID
	Status  OrderStatus `json:"status"`
	ForUser string      `json:"forUser"` // Ref to User.Username
}

// OrderStatus enum
type OrderStatus string

// This is a (sort of) enum of Order statuses
const (
	OrderNew        OrderStatus = "new"
	OrderReceived   OrderStatus = "received"
	OrderProcessing OrderStatus = "processing"
	OrderComplete   OrderStatus = "complete"
)

// OrderService defines core CRUD methods a orders service should have
type OrderService interface {
	GetOrder(orderID string) (*Order, error)
	GetOrdersForUser(username string) ([]string, error)
	ProcessOrder(order Order) error
	AddOrder(Order) error
	SetStatus(order *Order, status OrderStatus) error
	EmailNotify(Order) error
	SaveReport(Order) error
}

// Validate checks an order is correct
func Validate(o Order) error {
	if o.Amount <= 0 || len(o.Items) == 0 || o.Title == "" || o.ForUser == "" {
		return errors.New("Order failed validation")
	}
	return nil
}
