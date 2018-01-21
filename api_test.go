package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

var app App

func TestMain(m *testing.M) {
	app = App{}
	app.Initialize()
	code := m.Run()
	os.Exit(code)
}

func setupOrcs(count int) {
	// Empty orcStore
	orcStore = nil
	orcStore = make(map[string]Orc)
	id = 0

	if count == 0 {
		return
	}
	for i := 1; i <= count; i++ {
		orcStore[strconv.Itoa(i)] = Orc{"Urkhat", "Dabu", "DoomHammer", time.Now()}
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	app.Router.ServeHTTP(rec, req)

	return rec
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetOrcHandler(t *testing.T) {
	expected := 5
	setupOrcs(expected)

	// Test multiple orcs
	req, _ := http.NewRequest(http.MethodGet, "/api/orcs", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)

	// verify that we get expected number of orcs
	var orcs []Orc
	json.Unmarshal(res.Body.Bytes(), &orcs)
	actual := len(orcs)
	if actual != expected {
		t.Errorf("Expected length of %d. Got %d\n", expected, actual)
	}

	// Test single existing orc
	req, _ = http.NewRequest(http.MethodGet, "/api/orcs/1", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)

	var orc Orc
	json.Unmarshal(res.Body.Bytes(), &orc)
	checkContent(t, orc.Name, "Urkhat")

	// Test non-existent orc
	req, _ = http.NewRequest(http.MethodGet, "/api/orcs/666", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, res.Code)
}

func TestDeleteOrcHandler(t *testing.T) {
	setupOrcs(1)

	// Verify orc exists
	req, _ := http.NewRequest(http.MethodGet, "/api/orcs/1", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)

	// Delete orc
	req, _ = http.NewRequest(http.MethodDelete, "/api/orcs/1", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, res.Code)

	// Verify that orc doesn't exist
	req, _ = http.NewRequest(http.MethodGet, "/api/orcs/1", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, res.Code)

	// Delete non-existent orc
	req, _ = http.NewRequest(http.MethodDelete, "/api/orcs/666", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestPostOrcHandler(t *testing.T) {
	payload := []byte(`{"name": "Buldig", "greeting": "Swobu", "weapon":"DoomKitten"}`)

	req, _ := http.NewRequest(http.MethodPost, "/api/orcs", bytes.NewBuffer(payload))
	res := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, res.Code)
}

func TestPutOrcHandler(t *testing.T) {
	setupOrcs(1)
	key := 1

	// Verify orc exists
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/orcs/%d", key), nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)

	var orig map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &orig)

	// Update orc
	payload := []byte(`{"name": "Buldig", "greeting": "Swobu", "weapon":"DoomKitten"}`)
	req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/api/orcs/%d", key), bytes.NewBuffer(payload))
	res = executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, res.Code)

	// Get updated orc
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/orcs/%d", key), nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)

	var new map[string]interface{}
	var updated map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &updated)
	json.Unmarshal(payload, &new)

	if updated["name"] != new["name"] {
		t.Errorf("Expected orc name to change from %v to %v", orig["name"], new["name"])
	}

	// Update non-existent orc
	req, _ = http.NewRequest(http.MethodPut, "/api/orcs/666", bytes.NewBuffer(payload))
	res = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, res.Code)
}
