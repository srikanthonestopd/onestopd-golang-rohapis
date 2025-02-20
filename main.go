package main
Feature branch testing

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
	{OrderID: 3, CustomerName: "Jaswanth", OrderDate: time.Now().AddDate(0, 0, -4), TotalAmount: 900.00, Status: "Delivery on Tuesday"},
	{OrderID: 4, CustomerName: "Santosh", OrderDate: time.Now().AddDate(0, 0, -4), TotalAmount: 1900.00, Status: "Delivery on Wednesday"},
	{OrderID: 5, CustomerName: "Sravani", OrderDate: time.Now().AddDate(0, 0, -9), TotalAmount: 145.00, Status: "Delivery on Sunday"},
	{OrderID: 6, CustomerName: "yuvin", OrderDate: time.Now().AddDate(0, 0, -2), TotalAmount: 290.00, Status: "delivered"},
	{OrderID: 7, CustomerName: "vaishnavi", OrderDate: time.Now().AddDate(0, 0, -3), TotalAmount: 1000.00, Status: "to be shipped"},
	{OrderID: 8, CustomerName: "suraj", OrderDate: time.Now().AddDate(0, 0, -8), TotalAmount: 285.00, Status: " delivery by 6pm"},
	{OrderID: 9, CustomerName: "Sharath", OrderDate: time.Now().AddDate(0, 0, -6), TotalAmount: 185.00, Status: " delivery by 9pm"},
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
	log.Println("Server running on http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", router))
}
