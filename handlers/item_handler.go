package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"log"
	"net/http"
	"onestopd-golang-rohapis/config"
	"onestopd-golang-rohapis/models"
	"sort"
)

// Fetch all items, find the highest ID, and increment it
func getNextItemID() (string, error) {
	query := "SELECT id FROM `roh-apis`.`myscope`.`mycollection`;"
	rows, err := config.Cluster.Query(query, nil)
	if err != nil {
		return "", err
	}

	var ids []int
	for rows.Next() {
		var result map[string]interface{}
		err := rows.Row(&result)
		if err != nil {
			return "", err
		}

		// Convert ID to int (assuming numeric IDs)
		if idStr, ok := result["id"].(string); ok {
			var id int
			_, err := fmt.Sscanf(idStr, "item%d", &id)
			if err == nil {
				ids = append(ids, id)
			}
		}
	}

	// Sort and find the last ID
	sort.Ints(ids)
	newID := 1
	if len(ids) > 0 {
		newID = ids[len(ids)-1] + 1 // Increment last ID
	}

	return fmt.Sprintf("item%d", newID), nil
}

// Get the next available item ID
func GetNextItemIDHandler(w http.ResponseWriter, r *http.Request) {
	nextID, err := getNextItemID()
	if err != nil {
		http.Error(w, "Failed to get next item ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"next_id": nextID})
}

// Get data from Couchbase by ID
func GetData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// Fetch document from Couchbase
	result, err := config.Collection.Get(id, nil)
	if err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	var item models.Item
	err = result.Content(&item)
	if err != nil {
		http.Error(w, "Failed to parse data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Insert data into Couchbase with auto-incrementing ID
func InsertData(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Println("❌ JSON Decode Error:", err)
		return
	}

	// Generate the next unique ID
	item.ID, err = getNextItemID()
	if err != nil {
		http.Error(w, "Failed to generate unique ID", http.StatusInternalServerError)
		log.Println("❌ ID Generation Error:", err)
		return
	}

	// Print inserting data
	fmt.Printf("🔹 Inserting item with ID: %s\n", item.ID)

	// Insert into Couchbase
	_, err = config.Collection.Insert(item.ID, item, nil)
	if err != nil {
		http.Error(w, "Failed to insert data", http.StatusInternalServerError)
		log.Println("❌ Couchbase Insert Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data inserted successfully",
		"id":      item.ID,
	})
}
