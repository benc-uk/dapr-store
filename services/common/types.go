package common

// DaprState is a struct
type DaprState struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Product holds product data
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Cost        float32 `json:"cost"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
}

// Order holds information about a customer order
type Order struct {
	ID        string      `json:"id"`
	Amount    float32     `json:"amount"`
	ProductID string      `json:"productId"` // ref to a Product.ID
	Status    OrderStatus `json:"status"`
}

// OrderStatus enum
type OrderStatus string

// OrderNew blah blah
const (
	OrderNew        OrderStatus = "new"
	OrderReceived   OrderStatus = "received"
	OrderProcessing OrderStatus = "processing"
	OrderComplete   OrderStatus = "complete"
)
