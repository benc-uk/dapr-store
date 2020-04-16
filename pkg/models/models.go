package models

// Product holds product data
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Cost        float32 `json:"cost"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	OnOffer     bool    `json:"onOffer"`
}

// Order holds information about a customer order
type Order struct {
	ID      string      `json:"id"`
	Title   string      `json:"title"`
	Amount  float32     `json:"amount"`
	Items   []string    `json:"items"` // List of Product.ID
	Status  OrderStatus `json:"status"`
	ForUser string      `json:"forUser"` // Ref to User.Username
}

// User holds information about a registered user
type User struct {
	Username     string   `json:"username"`
	DisplayName  string   `json:"displayName"`
	ProfileImage string   `json:"profileImage"`
	Orders       []string `json:"orders"` // List of Order.IDs
}

// OrderStatus enum
type OrderStatus string

// This is a enum of Order statuses
const (
	OrderNew        OrderStatus = "new"
	OrderReceived   OrderStatus = "received"
	OrderProcessing OrderStatus = "processing"
	OrderComplete   OrderStatus = "complete"
)
