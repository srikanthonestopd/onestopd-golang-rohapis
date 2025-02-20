package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"onestopd-golang-rohapis/config"
	"onestopd-golang-rohapis/routes"
)

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	config.InitCouchbase() // Initialize Couchbase

	router := mux.NewRouter()
	routes.RegisterItemRoutes(router)
	routes.RegisterOrderRoutes(router)

	// Apply CORS middleware
	handler := enableCORS(router)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
	
}
