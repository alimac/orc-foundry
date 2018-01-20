package main

import "time"

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
