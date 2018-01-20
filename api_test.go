package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
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

func addOrcs(count int) {
	orcStore = nil
	orcStore = make(map[string]Orc)

	if count < 1 {
		count = 1
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

func TestGetOrcsHandler(t *testing.T) {
	addOrcs(5)
	req, _ := http.NewRequest(http.MethodGet, "/api/orcs", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
}

func TestGetOrcHandler(t *testing.T) {
	addOrcs(1)
	req, _ := http.NewRequest(http.MethodGet, "/api/orcs/1", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
}

func TestDeleteOrcHandler(t *testing.T) {
	addOrcs(1)

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
}

func TestPostOrcHandler(t *testing.T) {
	orcJSON := `{"name": "Buldig", "greeting": "Swobu", "weapon":"DoomKitten"}`

	req, _ := http.NewRequest(http.MethodPost, "/api/orcs", strings.NewReader(orcJSON))
	res := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, res.Code)
}

func TestPutOrcHandler(t *testing.T) {
	addOrcs(1)
	id := 1

	// Verify orc exists
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/orcs/%d", id), nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)

	var orig map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &orig)

	// Update orc
	payload := []byte(`{"name": "Buldig", "greeting": "Swobu", "weapon":"DoomKitten"}`)
	req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/api/orcs/%d", id), bytes.NewBuffer(payload))
	res = executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, res.Code)

	// Get updated orc
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/orcs/%d", id), nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)

	var new map[string]interface{}
	var updated map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &updated)
	json.Unmarshal(payload, &new)

	if updated["name"] != new["name"] {
		t.Errorf("Expected orc name to change from %v to %v", orig["name"], new["name"])
	}

}
