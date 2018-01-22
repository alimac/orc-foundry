package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// GetOrcHandler provides an endpoint for getting Orcs
func (a *App) GetOrcHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status, json := getItems(vars["id"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}

// PostOrcHandler provides an endpoint for creating new Orcs
func (a *App) PostOrcHandler(w http.ResponseWriter, r *http.Request) {
	var orc Orc
	// Decode the incoming Orc json
	err := json.NewDecoder(r.Body).Decode(&orc)
	if err != nil {
		panic(err)
	}

	orc.CreatedOn = time.Now()
	createItem(orc)

	json, err := json.Marshal(orc)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

// PutOrcHandler provides an endpoint for updating existing Orcs
func (a *App) PutOrcHandler(w http.ResponseWriter, r *http.Request) {
	var orc Orc

	vars := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&orc)
	status := updateItem(vars["id"], orc)

	w.WriteHeader(status)
}

// DeleteOrcHandler provides an endpoint for deleting existing Orcs
func (a *App) DeleteOrcHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := deleteItem(vars["id"])
	w.WriteHeader(status)
}
