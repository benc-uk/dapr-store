package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/benc-uk/dapr-store/pkg/problem"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprhttp "github.com/dapr/go-sdk/service/http"
	"github.com/gorilla/mux"
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
func NewService(serviceName string, router *mux.Router) *OrderService {
	storeName := env.GetEnvString("DAPR_STORE_NAME", "statestore")

	// Set up Dapr client & checks for Dapr sidecar, otherwise die
	client, err := dapr.NewClient()
	if err != nil {
		log.Panicln("FATAL! Dapr process/sidecar NOT found. Terminating!")
	}

	service := &OrderService{
		storeName,
		"orders-email",  // Hard coded, !TODO move to config env var
		"orders-report", // Hard coded, !TODO move to config env var
		serviceName,
		client,
	}

	// We don't actually use the service as we have one already
	// But we need to call AddTopicEventHandler to register the handler
	dummyService := daprhttp.NewServiceWithMux("notUsed", router)

	// Topic subscription
	var sub = &common.Subscription{
		PubsubName: env.GetEnvString("DAPR_PUBSUB_NAME", "pubsub"),
		Topic:      env.GetEnvString("DAPR_ORDERS_TOPIC", "orders-queue"),
		Route:      "/process-order",
	}

	if err := dummyService.AddTopicEventHandler(sub, service.pubSubOrderReceiver); err != nil {
		log.Printf("### Error adding topic subscription: %v", err)
	}
	// Suppress error, this will always fail, but we have to call Start to register the handler
	_ = dummyService.Start()

	return service
}

// AddOrder stores an order in Dapr state store
func (s *OrderService) AddOrder(order spec.Order) error {
	jsonPayload, err := json.Marshal(order)
	if err != nil {
		return problem.New500("err://json-marshall", "State JSON marshalling error", s.serviceName, nil, err)
	}
	if err := s.client.SaveState(context.Background(), s.storeName, order.ID, jsonPayload, nil); err != nil {
		return problem.NewDaprStateProblem(err, s.serviceName)
	}

	// This is a list of orderIDs for the user
	userOrders := []string{}
	// NOTE We use the username as a key in the orders state set, to hold an index of orders
	data, err := s.client.GetState(context.Background(), s.storeName, order.ForUser, nil)
	if err != nil {
		return problem.NewDaprStateProblem(err, s.serviceName)
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
		return problem.New500("err://json-marshall", "State JSON marshalling error", s.serviceName, nil, err)
	}
	if err := s.client.SaveState(context.Background(), s.storeName, order.ForUser, jsonPayload, nil); err != nil {
		log.Printf("### Error!, unable to save order list for user '%s'", order.ForUser)
		return problem.NewDaprStateProblem(err, s.serviceName)
	}

	return nil
}

// GetOrder fetches an order from Dapr state store
func (s *OrderService) GetOrder(orderID string) (*spec.Order, error) {
	data, err := s.client.GetState(context.Background(), s.storeName, orderID, nil)
	if err != nil {
		return nil, problem.NewDaprStateProblem(err, s.serviceName)
	}

	if data.Value == nil {
		prob := problem.New("err://not-found", "No data returned", 404, "Order: '"+orderID+"' not found", s.serviceName)
		return nil, prob
	}

	order := &spec.Order{}
	err = json.Unmarshal(data.Value, order)
	if err != nil {
		prob := problem.New("err://json-decode", "Malformed order JSON", 500, "JSON could not be decoded", s.serviceName)
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
		return errors.New("order not in correct status")
	}

	prob := s.AddOrder(order)
	if prob != nil {
		return prob
	}

	err = s.SetStatus(&order, spec.OrderReceived)
	if err != nil {
		log.Printf("### Failed to update state for order %s\n", err)
	}

	log.Printf("### Order %s was saved to state store\n", order.ID)

	// Save order to blob storage as a text file "report"
	// Also email to the user via SendGrid
	// For these to work configure the components in cmd/orders/components
	// If un-configured then nothing happens, and no output is send or generated

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

	// Fake background order processing
	time.AfterFunc(30*time.Second, func() {
		log.Printf("### Order %s is now processing\n", order.ID)
		_ = s.SetStatus(&order, spec.OrderProcessing)
	})

	// Fake background order completion
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
		return nil, problem.NewDaprStateProblem(err, s.serviceName)
	}

	orders := []string{}

	// If no data, just return an empty array
	if data.Value == nil {
		return orders, nil
	}

	err = json.Unmarshal(data.Value, &orders)
	if err != nil {
		prob := problem.New("err://json-decode", "Malformed orders JSON", 500, "JSON could not be decoded", s.serviceName)
		return nil, prob
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
		return problem.New500("err://json-marshall", "State JSON marshalling error", s.serviceName, nil, err)
	}
	if err := s.client.SaveState(context.Background(), s.storeName, order.ID, jsonPayload, nil); err != nil {
		log.Printf("### Error! Unable to update status of order '%s'", order.ID)
		return problem.NewDaprStateProblem(err, s.serviceName)
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

	err := s.client.InvokeOutputBinding(context.Background(), request)
	if err != nil {
		log.Printf("### Problem sending to email output: %+v", err)
		return problem.New500("dapr://binding", "Dapr binding invocation failed", s.serviceName, nil, err)
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

	err := s.client.InvokeOutputBinding(context.Background(), request)
	if err != nil {
		log.Printf("### Problem sending to blob / report output: %+v", err)
		return problem.New500("dapr://binding", "Dapr binding invocation failed", s.serviceName, nil, err)
	}

	return nil
}

// pubSubOrderReceiver is not part of the OrderService spec
// It is registered as the receiver for new messages on the Dapr pub/sub order topic
func (s *OrderService) pubSubOrderReceiver(ctx context.Context, e *common.TopicEvent) (bool, error) {
	order := spec.Order{}
	err := json.Unmarshal(e.RawData, &order)
	if err != nil {
		return false, err
	}

	err = s.ProcessOrder(order)
	if err != nil {
		return false, err
	}

	return false, nil
}
