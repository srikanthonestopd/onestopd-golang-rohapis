package models

import "time"

// Order struct for in-memory storage
type Order struct {
	OrderID      int       `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	OrderDate    time.Time `json:"order_date"`
	TotalAmount  float64   `json:"total_amount"`
	Status       string    `json:"status"`
}
