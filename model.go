package main

import (
	"encoding/json"
	"net/http"
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
