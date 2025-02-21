package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"onestopd-golang-rohapis/config"
	"onestopd-golang-rohapis/models"
	"strconv"
	"time"
)

// In-memory order list
var orders = []models.Order{
	{OrderID: 1, CustomerName: "John Doe", OrderDate: time.Now().AddDate(0, 0, -5), TotalAmount: 150.75, Status: "Delivered"},
	{OrderID: 2, CustomerName: "Jane Smith", OrderDate: time.Now().AddDate(0, 0, -3), TotalAmount: 200.00, Status: "Shipped"},
}

// Get all orders from Couchbase
func GetOrders(w http.ResponseWriter, r *http.Request) {
	query := "SELECT META().id, * FROM `roh-apis`.`myscope`.`orders`;"
	rows, err := config.Cluster.Query(query, nil)
	if err != nil {
		http.Error(w, "Failed to fetch orders from Couchbase", http.StatusInternalServerError)
		log.Println("❌ Couchbase Query Error:", err)
		return
	}

	var orders []models.Order
	for rows.Next() {
		var result map[string]interface{}
		err := rows.Row(&result)
		if err != nil {
			http.Error(w, "Failed to parse Couchbase data", http.StatusInternalServerError)
			return
		}

		// Convert result into an order struct
		var order models.Order
		orderMap := result["orders"].(map[string]interface{})
		order.OrderID = int(orderMap["order_id"].(float64))
		order.CustomerName = orderMap["customer_name"].(string)
		order.TotalAmount = orderMap["total_amount"].(float64)
		order.Status = orderMap["status"].(string) // ✅ Ensure status is included
		order.OrderDate, _ = time.Parse(time.RFC3339, orderMap["order_date"].(string))

		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// Get order by ID
func GetOrderById(w http.ResponseWriter, r *http.Request) {
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

// Add a new order to Couchbase
func AddOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate unique OrderID
	newOrder.OrderID = int(time.Now().Unix()) // Use timestamp as unique ID
	newOrder.OrderDate = time.Now()

	// ✅ Correct way to get collection from Couchbase
	ordersCollection := config.Cluster.Bucket("roh-apis").Scope("myscope").Collection("orders")

	// Insert into Couchbase
	_, err = ordersCollection.Insert(
		fmt.Sprintf("order_%d", newOrder.OrderID),
		newOrder,
		nil,
	)
	if err != nil {
		http.Error(w, "Failed to insert order into Couchbase", http.StatusInternalServerError)
		log.Println("Couchbase Insert Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Order added successfully",
		"order_id":  newOrder.OrderID,
		"orderDate": newOrder.OrderDate,
	})
}
