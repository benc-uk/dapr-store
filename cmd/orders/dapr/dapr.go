package dapr

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/dapr-store/pkg/dapr"
	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// OrderService is a Dapr based implementation of OrderService interface
type OrderService struct {
	*dapr.Helper
	storeName        string
	TopicName        string
	emailOutputName  string
	reportOutputName string
}

//
// New creates a new OrderService
//
func New(serviceName string) *OrderService {
	// Set up Dapr & checks for Dapr sidecar port, abort
	helper := dapr.NewHelper(serviceName)
	if helper == nil {
		os.Exit(1)
	}
	storeName := env.GetEnvString("DAPR_STORE_NAME", "statestore")
	topicName := env.GetEnvString("DAPR_ORDERS_TOPIC", "orders-queue")

	return &OrderService{
		helper,
		storeName,
		topicName,
		"orders-email",
		"orders-notify",
	}
}

func (d *OrderService) AddOrder(order spec.Order) error {
	// Call Dapr helper to save state
	prob := d.SaveState(d.storeName, order.ID, order)
	if prob != nil {
		return prob
	}

	userOrders := []string{}
	// !NOTE! We use the username as a key in the orders state set, to hold an index of orders
	data, prob := d.GetState(d.storeName, order.ForUser)
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
	prob = d.SaveState(d.storeName, order.ForUser, userOrders)
	if prob != nil {
		log.Printf("### Error!, unable to save order list for user '%s'", order.ForUser)
		return prob
	}

	return nil
}

func (d *OrderService) GetOrder(orderID string) (*spec.Order, error) {
	data, prob := d.GetState(d.storeName, orderID)
	if prob != nil {
		return nil, prob
	}

	if len(data) <= 0 {
		prob := problem.New("err://not-found", "No data returned", 404, "Order: '"+orderID+"' not found", d.ServiceName)
		return nil, prob
	}

	order := &spec.Order{}
	err := json.Unmarshal(data, order)
	if err != nil {
		prob := problem.New("err://json-decode", "Malformed order JSON", 500, "JSON could not be decoded", d.ServiceName)
		return nil, prob
	}

	return order, nil
}

func (d *OrderService) GetOrdersForUser(userName string) ([]string, error) {
	// !NOTE! We use the username as a key in the orders state set, to hold an index of orders
	data, prob := d.GetState(d.storeName, userName)
	if prob != nil {
		return nil, prob
	}

	orders := []string{}

	// If no data, just return an empty array
	if len(data) <= 0 {
		return orders, nil
	}

	log.Printf("\n}}}}}}} %s\n\n", string(data))
	err := json.Unmarshal(data, &orders)
	if err != nil {
		prob := problem.New("err://json-decode", "Malformed orders JSON", 500, "JSON could not be decoded", d.ServiceName)
		return nil, prob
	}

	return orders, nil
}

func (d *OrderService) SetStatus(order *spec.Order, status spec.OrderStatus) error {
	order.Status = status
	prob := d.SaveState(d.storeName, order.ID, order)
	if prob != nil {
		log.Printf("### Warning, order completion failed: %s", prob.Error())
		return prob
	}
	return nil
}

func (d *OrderService) EmailNotify(order spec.Order) error {
	emailMetadata := map[string]string{
		"emailTo": order.ForUser,
		"subject": "Dapr Store, order details: " + order.Title,
	}
	emailData := "<h1>Thanks for your order!</h1>Order title: " + order.Title + "<br>Order ID: " + order.ID +
		"<br>User: " + order.ForUser + "<br>Amount: £" + fmt.Sprintf("%.2f", order.Amount) + "<br><br>Enjoy your new dapper threads!"

	prob := d.SendOutput(d.emailOutputName, emailData, emailMetadata)
	if prob != nil {
		log.Printf("### Problem sending to email output: %+v", prob)
		return prob
	}

	return nil
}

func (d *OrderService) SaveReport(order spec.Order) error {
	blobMetadata := map[string]string{
		"ContentType": "text/plain",
		"blobName":    "order_" + order.ID + ".txt",
	}
	blobData := "----------\nOrder title:" + order.Title + "\nOrder ID: " + order.ID +
		"\nUser: " + order.ForUser + "\nAmount: " + fmt.Sprintf("%f", order.Amount) + "\n----------"

	prob := d.SendOutput(d.reportOutputName, blobData, blobMetadata)
	if prob != nil {
		log.Printf("### Problem sending to blob output: %+v", prob)
		return prob
	}

	return nil
}
