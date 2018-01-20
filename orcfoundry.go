package main

import (
	"html/template"
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
	templates["view"] = template.Must(template.ParseFiles("templates/view.html",
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
	viewModel = OrcModel{Orc{orc.Forge("name"), orc.Forge("greeting"),
		orc.Forge("weapon"), time.Now()}, "0"}

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
		renderTemplate(w, "edit", "base", viewModel)
	} else {
		http.Error(w, "Could not find the Orc to edit", http.StatusBadRequest)
	}
}

// editOrc
func getOrc(w http.ResponseWriter, r *http.Request) {
	var viewModel OrcModel

	// read value from route variable
	vars := mux.Vars(r)
	key := vars["id"]

	if orc, ok := orcStore[key]; ok {
		viewModel = OrcModel{orc, key}
		renderTemplate(w, "view", "base", viewModel)
	} else {
		http.Error(w, "Could not find the Orc to view", http.StatusNotFound)
	}

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
		http.Redirect(w, r, "/", 302)
	} else {
		http.Error(w, "Could not find the Orc to update", http.StatusBadRequest)
	}
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
		http.Redirect(w, r, "/", 302)
	} else {
		http.Error(w, "Could not find the Orc to delete", http.StatusBadRequest)
	}
}

func main() {
	app := App{}
	app.Initialize()
	app.Run(":8080")
}
