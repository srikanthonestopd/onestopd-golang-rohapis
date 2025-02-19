package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Order struct to represent order data
type Order struct {
	OrderID      int       `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	OrderDate    time.Time `json:"order_date"`
	TotalAmount  float64   `json:"total_amount"`
	Status       string    `json:"status"`
}

// In-memory database for demo purposes
var orders = []Order{
	{OrderID: 1, CustomerName: "John Doe", OrderDate: time.Now().AddDate(0, 0, -5), TotalAmount: 150.75, Status: "Delivered"},
	{OrderID: 2, CustomerName: "Jane Smith", OrderDate: time.Now().AddDate(0, 0, -3), TotalAmount: 200.00, Status: "Shipped"},
    {OrderID: 3, CustomerName: "Jaswanth", OrderDate: time.Now().AddDate(0, 0, -4),TotalAmount: 900.00, Status:"Delivery on Tuesday"},
}

// Get all orders
func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// Get order by ID
func getOrderById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	for _, order := range orders {
		if order.OrderID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(order)
			return
		}
	}

	http.Error(w, "Order not found", http.StatusNotFound)
}

// Add a new order
func addOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newOrder.OrderID = len(orders) + 1
	newOrder.OrderDate = time.Now()
	orders = append(orders, newOrder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
}

// Main function to start the server
func main() {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/api/orders", getOrders).Methods("GET")
	router.HandleFunc("/api/orders/{id}", getOrderById).Methods("GET")
	router.HandleFunc("/api/orders", addOrder).Methods("POST")

	// Start server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
