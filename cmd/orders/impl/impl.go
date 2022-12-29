package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/go-rest-api/pkg/dapr/pubsub"
	"github.com/benc-uk/go-rest-api/pkg/env"

	dapr "github.com/dapr/go-sdk/client"
)

// OrderService is a Dapr based implementation of OrderService interface
type OrderService struct {
	storeName        string
	emailOutputName  string
	reportOutputName string
	serviceName      string
	client           dapr.Client
}

// NewService creates a new OrderService
func NewService(serviceName string) *OrderService {
	storeName := env.GetEnvString("DAPR_STORE_NAME", "statestore")

	// Set up Dapr client & checks for Dapr sidecar, otherwise die
	client, err := dapr.NewClient()
	if err != nil {
		log.Panicln("FATAL! Dapr process/sidecar NOT found. Terminating!")
	}

	service := &OrderService{
		storeName,
		"orders-email",  // Hard coded, !TODO move to config env var\
		"orders-report", // Hard coded, !TODO move to config env var
		serviceName,
		client,
	}

	return service
}

// AddOrder stores an order in Dapr state store
func (s *OrderService) AddOrder(order spec.Order) error {
	jsonPayload, err := json.Marshal(order)
	if err != nil {
		return err
	}

	if err := s.client.SaveState(context.Background(), s.storeName, order.ID, jsonPayload, nil); err != nil {
		return err
	}

	// This is a list of orderIDs for the user
	userOrders := []string{}
	// NOTE We use the username as a key in the orders state set, to hold an index of orders
	data, err := s.client.GetState(context.Background(), s.storeName, order.ForUser, nil)
	if err != nil {
		return err
	}

	// Ignore any problem, it's possible it doesn't exist yet (user's first order)
	_ = json.Unmarshal(data.Value, &userOrders)

	alreadyExists := false

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

	// Save updated order list back, again keyed using user id
	jsonPayload, err = json.Marshal(userOrders)
	if err != nil {
		return err
	}

	if err := s.client.SaveState(context.Background(), s.storeName, order.ForUser, jsonPayload, nil); err != nil {
		log.Printf("### Error!, unable to save order list for user '%s'", order.ForUser)
		return err
	}

	return nil
}

// GetOrder fetches an order from Dapr state store
func (s *OrderService) GetOrder(orderID string) (*spec.Order, error) {
	data, err := s.client.GetState(context.Background(), s.storeName, orderID, nil)
	if err != nil {
		return nil, err
	}

	if data.Value == nil {
		return nil, OrderNotFoundError()
	}

	order := &spec.Order{}

	err = json.Unmarshal(data.Value, order)
	if err != nil {
		return nil, err
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
		return OrderStatusError()
	}

	if err := s.AddOrder(order); err != nil {
		return err
	}

	err = s.SetStatus(&order, spec.OrderReceived)
	if err != nil {
		log.Printf("### Failed to update state for order %s\n", err)
	}

	log.Printf("### Order %s was saved to state store\n", order.ID)

	// Save order to blob storage as a text file "report"
	// Also email to the user via SendGrid
	// For these to work configure the components in cmd/orders/components
	// If un-configured then nothing happens (maybe some errors are logged)

	// Currently the SendGrid integration in Dapr is fubar
	// To be fixed by this PR https://github.com/dapr/components-contrib/pull/1867
	err = s.EmailNotify(order)
	if err != nil {
		log.Printf("### Email notification failed %s\n", err)
	}

	err = s.SaveReport(order)
	if err != nil {
		log.Printf("### Saving order report failed %s\n", err)
	}

	// Fake background order processing, move to processing after 30 seconds
	time.AfterFunc(30*time.Second, func() {
		log.Printf("### Order %s is now processing\n", order.ID)
		_ = s.SetStatus(&order, spec.OrderProcessing)
	})

	// Fake background order completion, move to complete after 2 minutes
	time.AfterFunc(120*time.Second, func() {
		log.Printf("### Order %s completed\n", order.ID)
		_ = s.SetStatus(&order, spec.OrderComplete)
	})

	return nil
}

// GetOrdersForUser fetches a list of order ids for a given username
func (s *OrderService) GetOrdersForUser(userName string) ([]string, error) {
	// NOTE We use the username as a key in the orders state set, to hold an index of orders
	data, err := s.client.GetState(context.Background(), s.storeName, userName, nil)
	if err != nil {
		return nil, err
	}

	orders := []string{}

	// If no data, just return an empty array
	if data.Value == nil {
		return orders, nil
	}

	if err = json.Unmarshal(data.Value, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// SetStatus updates the status of an order
func (s *OrderService) SetStatus(order *spec.Order, status spec.OrderStatus) error {
	log.Printf("### Setting status for order %s to %s\n", order.ID, status)
	order.Status = status

	// Save updated order list back, again keyed using user id
	jsonPayload, err := json.Marshal(order)
	if err != nil {
		return err
	}

	if err := s.client.SaveState(context.Background(), s.storeName, order.ID, jsonPayload, nil); err != nil {
		log.Printf("### Error! Unable to update status of order '%s'", order.ID)
		return err
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

	request := &dapr.InvokeBindingRequest{}
	request.Metadata = emailMetadata
	request.Data = []byte(emailData)
	request.Name = s.emailOutputName
	request.Operation = "create"

	if err := s.client.InvokeOutputBinding(context.Background(), request); err != nil {
		log.Printf("### Problem sending to email output: %s", err.Error())
		return err
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

	request := &dapr.InvokeBindingRequest{}
	request.Metadata = blobMetadata
	request.Data = []byte(blobData)
	request.Name = s.reportOutputName
	request.Operation = "create"

	if err := s.client.InvokeOutputBinding(context.Background(), request); err != nil {
		log.Printf("### Problem sending to blob / report output: %s", err.Error())
		return err
	}

	return nil
}

// pubSubOrderReceiver is an adaptor of sorts, not part of the OrderService spec
// It is registered as the receiver for new messages on the Dapr pub/sub order topic
func (s *OrderService) PubSubOrderReceiver(event *pubsub.CloudEvent) error {
	// This JSON nonsense is an "easy" way to convert
	// The event.Data which is a map back into a real Order
	jsonData, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	var order spec.Order
	if err := json.Unmarshal(jsonData, &order); err != nil {
		return err
	}

	// Now we have a real order, we can process it
	if err := s.ProcessOrder(order); err != nil {
		return err
	}

	return nil
}
