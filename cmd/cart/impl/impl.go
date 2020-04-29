package impl

import (
	"math/rand"
	"os"
	"time"

	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/dapr-store/pkg/dapr"
	"github.com/benc-uk/dapr-store/pkg/env"
)

type CartService struct {
	*dapr.Helper
	topicName string
}

//
// New creates a new OrderService
//
func NewService(serviceName string) *CartService {
	// Set up Dapr & checks for Dapr sidecar port, abort
	helper := dapr.NewHelper(serviceName)
	if helper == nil {
		os.Exit(1)
	}
	topicName := env.GetEnvString("DAPR_ORDERS_TOPIC", "orders-queue")

	return &CartService{
		helper,
		topicName,
	}
}

func (s CartService) SubmitOrder(order *spec.Order) error {
	order.ID = makeID(5)
	order.Status = spec.OrderNew

	prob := s.PublishMessage(s.topicName, order)
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
