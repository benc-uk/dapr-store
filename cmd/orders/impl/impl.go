package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/dapr-store/pkg/dapr"
	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/benc-uk/dapr-store/pkg/problem"
	"github.com/gorilla/mux"
)

// OrderService is a Dapr based implementation of OrderService interface
type OrderService struct {
	*dapr.Helper
	storeName        string
	emailOutputName  string
	reportOutputName string
}

// NewService creates a new OrderService
func NewService(serviceName string, router *mux.Router) *OrderService {
	// Set up Dapr & checks for Dapr sidecar port, abort
	helper := dapr.NewHelper(serviceName)
	if helper == nil {
		os.Exit(1)
	}
	storeName := env.GetEnvString("DAPR_STORE_NAME", "statestore")

	service := &OrderService{
		helper,
		storeName,
		"orders-email",
		"orders-report",
	}

	// Dapr pub/sub specific
	topicName := env.GetEnvString("DAPR_ORDERS_TOPIC", "orders-queue")
	helper.RegisterTopicSubscriptions([]string{topicName}, router)
	helper.RegisterTopicReceiver(topicName, router, service.pubSubOrderReceiver)

	return service
}

// AddOrder stores an order in Dapr state store
func (s *OrderService) AddOrder(order spec.Order) error {
	// Call Dapr helper to save state
	prob := s.SaveState(s.storeName, order.ID, order)
	if prob != nil {
		return prob
	}

	userOrders := []string{}
	// !NOTE! We use the username as a key in the orders state set, to hold an index of orders
	data, prob := s.GetState(s.storeName, order.ForUser)
	// Ignore any problem, it's possible it doesn't exist yet (user's first order)
	_ = json.Unmarshal(data, &userOrders)

	alreadyExists := false
	log.Printf("### userOrders is %v", userOrders)
	for _, oid := range userOrders {
		if order.ID == oid {
			alreadyExists = true
		}
	}

	if !alreadyExists {
		userOrders = append(userOrders, order.ID)
	} else {
		log.Printf("### Warning, duplicate order '%s' for user '%s' detected", order.ID, order.ForUser)
	}

	// Save new list back
	prob = s.SaveState(s.storeName, order.ForUser, userOrders)
	if prob != nil {
		log.Printf("### Error!, unable to save order list for user '%s'", order.ForUser)
		return prob
	}

	return nil
}

// GetOrder fetches an order from Dapr state store
func (s *OrderService) GetOrder(orderID string) (*spec.Order, error) {
	data, prob := s.GetState(s.storeName, orderID)
	if prob != nil {
		return nil, prob
	}

	if len(data) <= 0 {
		prob := problem.New("err://not-found", "No data returned", 404, "Order: '"+orderID+"' not found", s.ServiceName)
		return nil, prob
	}

	order := &spec.Order{}
	err := json.Unmarshal(data, order)
	if err != nil {
		prob := problem.New("err://json-decode", "Malformed order JSON", 500, "JSON could not be decoded", s.ServiceName)
		return nil, prob
	}

	return order, nil
}

// ProcessOrder is a fake set of order processing tasks, mainly progressing the status
func (s *OrderService) ProcessOrder(order spec.Order) error {
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
	time.AfterFunc(30*time.Second, func() {
		s.SetStatus(&order, spec.OrderProcessing)
	})

	// Fake background order completion
	time.AfterFunc(120*time.Second, func() {
		s.SetStatus(&order, spec.OrderComplete)
	})

	return nil
}

// GetOrdersForUser fetches a list of order ids for a given username
func (s *OrderService) GetOrdersForUser(userName string) ([]string, error) {
	// !NOTE! We use the username as a key in the orders state set, to hold an index of orders
	data, prob := s.GetState(s.storeName, userName)
	if prob != nil {
		return nil, prob
	}

	orders := []string{}

	// If no data, just return an empty array
	if len(data) <= 0 {
		return orders, nil
	}

	err := json.Unmarshal(data, &orders)
	if err != nil {
		prob := problem.New("err://json-decode", "Malformed orders JSON", 500, "JSON could not be decoded", s.ServiceName)
		return nil, prob
	}

	return orders, nil
}

// SetStatus updates the status of an order
func (s *OrderService) SetStatus(order *spec.Order, status spec.OrderStatus) error {
	order.Status = status
	prob := s.SaveState(s.storeName, order.ID, order)
	if prob != nil {
		log.Printf("### Warning, order completion failed: %s", prob.Error())
		return prob
	}
	return nil
}

// EmailNotify uses Dapr SendGrid output binding to send an email
func (s *OrderService) EmailNotify(order spec.Order) error {
	emailMetadata := map[string]string{
		"emailTo": order.ForUser,
		"subject": "Dapr Store, order details: " + order.Title,
	}
	emailData := "<h1>Thanks for your order!</h1>Order title: " + order.Title + "<br>Order ID: " + order.ID +
		"<br>User: " + order.ForUser + "<br>Amount: £" + fmt.Sprintf("%.2f", order.Amount) + "<br><br>Enjoy your new dapper threads!"

	prob := s.SendOutput(s.emailOutputName, emailData, emailMetadata)
	if prob != nil {
		log.Printf("### Problem sending to email output: %+v", prob)
		return prob
	}

	return nil
}

// SaveReport uses Dapr Azure Blob output binding to store a order report
func (s *OrderService) SaveReport(order spec.Order) error {
	blobName := "order_" + order.ID + ".txt"
	blobMetadata := map[string]string{
		"ContentType": "text/plain",
		"blobName":    blobName,
	}
	blobData := "----------\nTitle: " + order.Title + "\nOrder ID: " + order.ID +
		"\nUser: " + order.ForUser + "\nAmount: " + fmt.Sprintf("%f", order.Amount) + "\n----------"

	log.Printf("### Saving report to blob: %s", blobName)
	prob := s.SendOutput(s.reportOutputName, blobData, blobMetadata)
	if prob != nil {
		log.Printf("### Problem sending to blob output: %+v", prob)
		return prob
	}

	return nil
}

// pubSubOrderReceiver is an adaptor of sorts, not part of the OrderService spec
// It is registered as the receiver for new messages on the Dapr pub/sub order topic
func (s *OrderService) pubSubOrderReceiver(data io.Reader) error {
	type cloudevent struct {
		ID   string     `json:"id"`
		Data spec.Order `json:"data"`
	}
	event := &cloudevent{}

	err := json.NewDecoder(data).Decode(&event)
	if err != nil {
		return err
	}

	order := event.Data
	err = s.ProcessOrder(order)
	if err != nil {
		return err
	}

	return nil
}
