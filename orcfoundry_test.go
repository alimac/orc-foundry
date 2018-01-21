package main

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func checkContent(t *testing.T, content string, expected string) {
	if !strings.Contains(content, expected) {
		t.Logf("Content: %q", content)
		t.Errorf("String '%s' missing from body: %q", expected, content)
	}
}

func TestApp(t *testing.T) {
	a := App{}
	a.Initialize()

	// Use a goroutine to run the app to serve requests and exit
	go func() {
		defer a.Server.Close()
		a.Run(":3000")
	}()

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:3000/", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
	checkContent(t, res.Body.String(), "Orcs!")
}

func TestGetOrcs(t *testing.T) {
	setupOrcs(1)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
	checkContent(t, res.Body.String(), "Holding DoomHammer")
}

func TestEditOrc(t *testing.T) {
	setupOrcs(1)

	// Edit an existing orc
	req, _ := http.NewRequest(http.MethodGet, "/orcs/edit/1", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
	checkContent(t, res.Body.String(), "Reforge an orc")

	// Edit non-existent orc
	req, _ = http.NewRequest(http.MethodGet, "/orcs/edit/666", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestDeleteOrc(t *testing.T) {
	setupOrcs(1)

	// Delete existing orc
	req, _ := http.NewRequest(http.MethodDelete, "/orcs/delete/1", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusFound, res.Code)

	// Verify orc does not exist
	req, _ = http.NewRequest(http.MethodGet, "/orcs/view/1", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, res.Code)

	// Delete non-existent orc
	req, _ = http.NewRequest(http.MethodDelete, "/orcs/delete/666", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, res.Code)
}

func TestCreateOrc(t *testing.T) {
	payload := url.Values{"name": {"Gonmund"}, "greeting": {"Fubu"}, "weapon": {"AgonySickle"}}

	// Add an orc
	req, _ := http.NewRequest(http.MethodGet, "/orcs/add", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
	checkContent(t, res.Body.String(), "Forge a new orc")

	// Save an orc
	req, _ = http.NewRequest(http.MethodPost, "/orcs/save", bytes.NewBufferString(payload.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res = executeRequest(req)
	checkResponseCode(t, http.StatusFound, res.Code)

	// Verify orc exists
	req, _ = http.NewRequest(http.MethodGet, "/orcs/view/1", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
	checkContent(t, res.Body.String(), "Meet Gonmund")
}

func TestUpdateOrc(t *testing.T) {
	setupOrcs(1)
	payload := url.Values{"name": {"Gonmund"}, "greeting": {"Fubu"}, "weapon": {"AgonySickle"}}

	// Update orc
	req, _ := http.NewRequest(http.MethodPut, "/orcs/update/1", bytes.NewBufferString(payload.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := executeRequest(req)
	checkResponseCode(t, http.StatusFound, res.Code)

	// Update non-existent orc
	req, _ = http.NewRequest(http.MethodPut, "/orcs/update/666", bytes.NewBufferString(payload.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, res.Code)

	// Verify orc is updated
	req, _ = http.NewRequest(http.MethodGet, "/orcs/view/1", nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
	checkContent(t, res.Body.String(), "AgonySickle")
}
