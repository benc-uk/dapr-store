package main

// Order is a struct
type Order struct {
	ID        string  `json:"id"`
	Amount    float32 `json:"amount"`
	ProductID string  `json:"productId"`
}
