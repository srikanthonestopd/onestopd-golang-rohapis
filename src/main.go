package main

import (
	"encoding/json"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Define Couchbase cluster connection
var cluster *gocb.Cluster
var collection *gocb.Collection

// Connect to Couchbase
func initCouchbase() {
	var err error
	cluster, err = gocb.Connect("couchbase://localhost", gocb.ClusterOptions{
		Username: "Administrator",
		Password: "admin123",
	})
	if err != nil {
		log.Fatalf("Failed to connect to Couchbase: %v", err)
	}

	// Open bucket and collection
	bucket := cluster.Bucket("roh-api")
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatalf("Bucket is not ready: %v", err)
	}

	collection = bucket.Scope("myscope").Collection("mycollection")
	fmt.Println("✅ Connected to Couchbase and collection is ready!")
}

// Define Data Model
type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// POST API to Insert Data
func insertData(w http.ResponseWriter, r *http.Request) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = collection.Insert(item.ID, item, nil)
	if err != nil {
		http.Error(w, "Failed to insert data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Data inserted successfully"})
}

// GET API to Fetch Data
func getData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	result, err := collection.Get(id, nil)
	if err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	var item Item
	err = result.Content(&item)
	if err != nil {
		http.Error(w, "Failed to parse data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

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
	{OrderID: 5, CustomerName: "Sravani", OrderDate: time.Now().AddDate(0, 0, -4), TotalAmount: 145.00, Status: "Delivery on Sunday"},
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
	initCouchbase()

	router := mux.NewRouter()
	router.HandleFunc("/api/insert", insertData).Methods("POST")
	router.HandleFunc("/api/get/{id}", getData).Methods("GET")

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
