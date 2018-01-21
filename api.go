package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// GetOrcHandler provides an endpoint for getting Orcs
func GetOrcHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	if key != "" {
		// Retrieve from store
		if orc, ok := orcStore[key]; ok {
			w.Header().Set("Content-Type", "application/json")

			json, err := json.Marshal(orc)
			if err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(json)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		var orcs []Orc

		for _, v := range orcStore {
			orcs = append(orcs, v)
		}
		w.Header().Set("Content-Type", "application/json")
		json, err := json.Marshal(orcs)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

// PostOrcHandler provides an endpoint for creating new Orcs
func PostOrcHandler(w http.ResponseWriter, r *http.Request) {

	var orc Orc
	// Decode the incoming Orc json
	err := json.NewDecoder(r.Body).Decode(&orc)
	if err != nil {
		panic(err)
	}

	orc.CreatedOn = time.Now()
	id++
	key := strconv.Itoa(id)
	orcStore[key] = orc

	json, err := json.Marshal(orc)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)

}

// PutOrcHandler provides an endpoint for updating existing Orcs
func PutOrcHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	key := vars["id"]
	var orcToUpdate Orc

	// Decode the incoming Orc json
	err = json.NewDecoder(r.Body).Decode(&orcToUpdate)
	if err != nil {
		panic(err)
	}

	if orc, ok := orcStore[key]; ok {
		orcToUpdate.CreatedOn = orc.CreatedOn
		// delete the existing item and add the updated item
		delete(orcStore, key)
		orcStore[key] = orcToUpdate
		w.WriteHeader(http.StatusNoContent)
	} else {
		log.Printf("Could not find key of Orc %s to update", key)
		w.WriteHeader(http.StatusBadRequest)
	}

}

// DeleteOrcHandler provides an endpoint for deleting existing Orcs
func DeleteOrcHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	// Remove from store
	if _, ok := orcStore[key]; ok {
		// Delete existing item
		delete(orcStore, key)
		w.WriteHeader(http.StatusNoContent)
	} else {
		log.Printf("Could not find key of Orc %s to delete", key)
		w.WriteHeader(http.StatusBadRequest)
	}
}
