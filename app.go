package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// App struct contains app config
type App struct {
	Router *mux.Router
}

// Initialize sets up the app configuration and initializes
// app routes
func (app *App) Initialize() {
	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

// Run runs the app
func (app *App) Run(port string) {
	host := "127.0.0.1"
	server := &http.Server{
		Handler:      app.Router,
		Addr:         fmt.Sprintf("%s%s", host, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/", getOrcs)
	app.Router.HandleFunc("/orcs/", getOrcs)
	app.Router.HandleFunc("/orcs/view/{id}", getOrc)
	app.Router.HandleFunc("/orcs/add", addOrc)
	app.Router.HandleFunc("/orcs/save", saveOrc)
	app.Router.HandleFunc("/orcs/edit/{id}", editOrc)
	app.Router.HandleFunc("/orcs/update/{id}", updateOrc)
	app.Router.HandleFunc("/orcs/delete/{id}", deleteOrc)

	app.Router.HandleFunc("/api/orcs", GetOrcsHandler).Methods("GET")
	app.Router.HandleFunc("/api/orcs/{id}", GetOrcHandler).Methods("GET")
	app.Router.HandleFunc("/api/orcs", PostOrcHandler).Methods("POST")
	app.Router.HandleFunc("/api/orcs/{id}", PutOrcHandler).Methods("PUT")
	app.Router.HandleFunc("/api/orcs/{id}", DeleteOrcHandler).Methods("DELETE")
}
