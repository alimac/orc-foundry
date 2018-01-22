package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alimac/orc"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// App struct contains app config
type App struct {
	Router *mux.Router
	Server *http.Server
	Port   string
	Host   string
}

// Initialize sets up the app configuration and initializes
// app routes
func (app *App) Initialize() {
	app.Router = mux.NewRouter()
	app.initializeRoutes()
	app.initializeOrcs()

	// Heroku - get port from environment
	port := os.Getenv("PORT")

	// Local and CI - set host and port
	if port == "" {
		app.Port = ":8080"
		app.Host = "127.0.0.1"
	} else {
		app.Port = port
		app.Host = ":"
	}
}

// Run runs the app
func (app *App) Run() {
	loggedRouter := handlers.LoggingHandler(os.Stdout, app.Router)
	app.Server = &http.Server{
		Handler:      loggedRouter,
		Addr:         fmt.Sprintf("%s%s", app.Host, app.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(app.Server.ListenAndServe())
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/", app.getOrc)
	app.Router.HandleFunc("/orcs/view/{id}", app.getOrc)
	app.Router.HandleFunc("/orcs/add", app.addOrc)
	app.Router.HandleFunc("/orcs/save", app.saveOrc)
	app.Router.HandleFunc("/orcs/edit/{id}", app.editOrc)
	app.Router.HandleFunc("/orcs/update/{id}", app.updateOrc)
	app.Router.HandleFunc("/orcs/delete/{id}", app.deleteOrc)

	app.Router.HandleFunc("/api/orcs", app.GetOrcHandler).Methods("GET")
	app.Router.HandleFunc("/api/orcs/{id}", app.GetOrcHandler).Methods("GET")
	app.Router.HandleFunc("/api/orcs", app.PostOrcHandler).Methods("POST")
	app.Router.HandleFunc("/api/orcs/{id}", app.PutOrcHandler).Methods("PUT")
	app.Router.HandleFunc("/api/orcs/{id}", app.DeleteOrcHandler).Methods("DELETE")
}

func (app *App) initializeOrcs() {
	// Initalize app with 5 orcs
	for i := 0; i < 5; i++ {
		createItem(Orc{
			orc.Forge("name"),
			orc.Forge("greeting"),
			orc.Forge("weapon"),
			time.Now(),
		})
	}
}
