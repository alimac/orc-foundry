package main

import (
  "encoding/json"
	//"fmt"
	"log"
	"net/http"
  "strconv"
	"time"

  "github.com/gorilla/mux"
)

type Orc struct {
  Name      string    `json:"name"`
  Greeting  string    `json:"greeting"`
  CreatedOn time.Time `json:"createdon"`
}

// Store for the Orcs collection
var orcStore = make(map[string]Orc)

// Variable to generate key for the collection
var id int = 0

// HTTP Get - /api/orcs
func GetOrcHandler(w http.ResponseWriter, r *http.Request) {
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

// HTTP Post - /api/orcs
func PostOrcHandler(w http.ResponseWriter, r *http.Request) {

  var orc Orc
  // Decode the incoming Orc json
  err := json.NewDecoder(r.Body).Decode(&orc)
  if err != nil {
    panic(err)
  }

  orc.CreatedOn = time.Now()
  id++
  k := strconv.Itoa(id)
  orcStore[k] = orc

  json, err := json.Marshal(orc)
  if err != nil {
    panic(err)
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  w.Write(json)

}

func PutOrcHandler(w http.ResponseWriter, r *http.Request) {
  var err error
  vars := mux.Vars(r)
  k := vars["id"]
  var orcToUpdate Orc

  // Decode the incoming Orc json
  err = json.NewDecoder(r.Body).Decode(&orcToUpdate)
  if err != nil {
    panic(err)
  }

  if orc, ok := orcStore[k]; ok  {
    orcToUpdate.CreatedOn = orc.CreatedOn
    // delete the existing item and add the updated item
    delete(orcStore, k)
    orcStore[k] = orcToUpdate
  } else {
    log.Printf("Could not find key of Orc %s to update", k)
  }
  w.WriteHeader(http.StatusNoContent)
}

func DeleteOrcHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  k := vars["id"]
  // Remove from store
  if _, ok := orcStore[k]; ok {
    // Delete existing item
    delete(orcStore, k)
  } else {
    log.Printf("Could not find key of Orc %s to delete", k)
  }
  w.WriteHeader(http.StatusNoContent)
}


func main() {

	r := mux.NewRouter().StrictSlash(false)
  r.HandleFunc("/api/orcs", GetOrcHandler).Methods("GET")
  r.HandleFunc("/api/orcs", PostOrcHandler).Methods("POST")
  r.HandleFunc("/api/orcs/{id}", PutOrcHandler).Methods("PUT")
  r.HandleFunc("/api/orcs/{id}", DeleteOrcHandler).Methods("DELETE")

  server := &http.Server {
    Addr: ":8080",
    Handler: r,
  }

	server.ListenAndServe()
}

