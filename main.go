package main

import (
	"encoding/json"
	"html/template"
	//"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alimac/orc"
	"github.com/gorilla/mux"
)

var templates map[string]*template.Template

// Compile view templates
func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(template.ParseFiles("templates/index.html",
		"templates/base.html"))
	templates["add"] = template.Must(template.ParseFiles("templates/add.html",
		"templates/base.html"))
	templates["edit"] = template.Must(template.ParseFiles("templates/edit.html",
		"templates/base.html"))
}

// Render templates for the given name, template definition and data object
func renderTemplate(w http.ResponseWriter, name string, template string, viewModel interface{}) {
	// Ensure template exists in the map
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "Template does not exist", http.StatusInternalServerError)
	}
	err := tmpl.ExecuteTemplate(w, template, viewModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getOrcs
func getOrcs(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", "base", orcStore)
}

// addOrc
func addOrc(w http.ResponseWriter, r *http.Request) {
	var viewModel OrcModel
	viewModel = OrcModel{Orc{orc.GenerateName(), orc.GenerateGreeting(),
		orc.GenerateWeapon(), time.Now()}, "0"}

	renderTemplate(w, "add", "base", viewModel)
}

// saveOrc
func saveOrc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.PostFormValue("name")
	greeting := r.PostFormValue("greeting")
	weapon := r.PostFormValue("weapon")
	orc := Orc{name, greeting, weapon, time.Now()}

	// increment value of id for generating key for the map
	id++
	// convert id value to string
	key := strconv.Itoa(id)
	orcStore[key] = orc
	http.Redirect(w, r, "/", 302)
}

// editOrc
func editOrc(w http.ResponseWriter, r *http.Request) {
	var viewModel OrcModel

	// read value from route variable
	vars := mux.Vars(r)
	key := vars["id"]

	if orc, ok := orcStore[key]; ok {
		viewModel = OrcModel{orc, key}
	} else {
		http.Error(w, "Could not find the Orc to edit", http.StatusBadRequest)
	}
	renderTemplate(w, "edit", "base", viewModel)
}

// updateOrc
func updateOrc(w http.ResponseWriter, r *http.Request) {
	// Read values from route variable
	vars := mux.Vars(r)
	key := vars["id"]
	var orcToUpdate Orc
	if orc, ok := orcStore[key]; ok {
		r.ParseForm()
		orcToUpdate.Name = r.PostFormValue("name")
		orcToUpdate.Greeting = r.PostFormValue("greeting")
		orcToUpdate.Weapon = r.PostFormValue("weapon")
		orcToUpdate.CreatedOn = orc.CreatedOn

		// delete existing item and add the updated item
		delete(orcStore, key)
		orcStore[key] = orcToUpdate

	} else {
		http.Error(w, "Could not find the Orc to update", http.StatusBadRequest)
	}
	http.Redirect(w, r, "/", 302)
}

// deleteOrc is a handler for "/orcs/delete/{id}" which deletes an item from the store
func deleteOrc(w http.ResponseWriter, r *http.Request) {
	// read value from the route Variable
	vars := mux.Vars(r)
	key := vars["id"]
	// Remove from the Store
	if _, ok := orcStore[key]; ok {
		// delete existing item
		delete(orcStore, key)
	} else {
		http.Error(w, "Could not find the Orc to delete", http.StatusBadRequest)
	}
	http.Redirect(w, r, "/", 302)
}

// Orc type provides an orc with a name and greeting
type Orc struct {
	Name      string    `json:"name"`
	Greeting  string    `json:"greeting"`
	Weapon    string    `json:"weapon"`
	CreatedOn time.Time `json:"createdon"`
}

// OrcModel is a view model for editing Orcs
type OrcModel struct {
	Orc
	ID string
}

// Store for the Orcs collection
var orcStore = make(map[string]Orc)

// Variable to generate key for the collection
var id int

// GetOrcHandler provides an endpoint for getting Orcs
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
	} else {
		log.Printf("Could not find key of Orc %s to update", key)
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteOrcHandler provides an endpoint for deleting existing Orcs
func DeleteOrcHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	// Remove from store
	if _, ok := orcStore[key]; ok {
		// Delete existing item
		delete(orcStore, key)
	} else {
		log.Printf("Could not find key of Orc %s to delete", key)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {

	// Create sample orcs
	id++
	orcStore[strconv.Itoa(id)] = Orc{"Urkhat", "Dabu", "DoomHammer", time.Now()}
	id++
	orcStore[strconv.Itoa(id)] = Orc{"Pigdug", "Zub zub", "DeathKettle", time.Now()}


	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", getOrcs)
	r.HandleFunc("/orcs/", getOrcs)
	r.HandleFunc("/orcs/add", addOrc)
	r.HandleFunc("/orcs/save", saveOrc)
	r.HandleFunc("/orcs/edit/{id}", editOrc)
	r.HandleFunc("/orcs/update/{id}", updateOrc)
	r.HandleFunc("/orcs/delete/{id}", deleteOrc)

	r.HandleFunc("/api/orcs", GetOrcHandler).Methods("GET")
	r.HandleFunc("/api/orcs", PostOrcHandler).Methods("POST")
	r.HandleFunc("/api/orcs/{id}", PutOrcHandler).Methods("PUT")
	r.HandleFunc("/api/orcs/{id}", DeleteOrcHandler).Methods("DELETE")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}
