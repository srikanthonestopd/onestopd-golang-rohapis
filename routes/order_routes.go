package routes

import (
	"github.com/gorilla/mux"
	"onestopd-golang-rohapis/handlers"
)

func RegisterOrderRoutes(router *mux.Router) {
	router.HandleFunc("/api/orders", handlers.GetOrders).Methods("GET")
	router.HandleFunc("/api/orders", handlers.AddOrder).Methods("POST")
}
