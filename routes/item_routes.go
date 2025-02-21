package routes

import (
	"github.com/gorilla/mux"
	"onestopd-golang-rohapis/handlers"
)

func RegisterItemRoutes(router *mux.Router) {
	router.HandleFunc("/api/insert", handlers.InsertData).Methods("POST")
	router.HandleFunc("/api/get/{id}", handlers.GetData).Methods("GET")
	router.HandleFunc("/api/next-id", handlers.GetNextItemIDHandler).Methods("GET")

}
