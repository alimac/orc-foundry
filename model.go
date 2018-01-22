package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Orc type provides an orc with a name and greeting
type Orc struct {
	Name      string    `json:"name"`
	Greeting  string    `json:"greeting"`
	Weapon    string    `json:"weapon"`
	CreatedOn time.Time `json:"createdon"`
}

// Store for the Orcs collection
var orcStore = make(map[string]Orc)

// Variable to generate key for the collection
var id int

// OrcModel is a view model for editing Orcs
type OrcModel struct {
	Orc
	ID string
}

func getItems(key string) (status int, items []byte) {
	// Retrieve specific item
	if key != "" {
		if orc, ok := orcStore[key]; ok {
			items, _ := json.Marshal(orc)
			return http.StatusOK, items
		}
		return http.StatusNotFound, nil
	}
	// Retrieve all items
	var orcs []Orc
	for _, v := range orcStore {
		orcs = append(orcs, v)
	}
	items, _ = json.Marshal(orcs)
	return http.StatusOK, items
}

func createItem(orc Orc) {
	// Increment id
	id++
	// Convert id value to string
	key := strconv.Itoa(id)
	// Store item
	orcStore[key] = orc
}

func updateItem(key string, o Orc) (status int) {
	if orc, ok := orcStore[key]; ok {
		// Retain created on timestamp
		o.CreatedOn = orc.CreatedOn
		// delete the existing item and add the updated item
		delete(orcStore, key)
		orcStore[key] = o
		return http.StatusNoContent
	}
	log.Printf("Could not find key of Orc %s to update", key)
	return http.StatusBadRequest
}

func deleteItem(key string) (status int) {
	if _, ok := orcStore[key]; ok {
		// Delete existing item
		delete(orcStore, key)
		return http.StatusNoContent
	}
	log.Printf("Could not find key of Orc %s to delete", key)
	return http.StatusBadRequest
}
