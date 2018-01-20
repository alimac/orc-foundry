package main

import (
	"bytes"
	"fmt"
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

func TestGetOrcs(t *testing.T) {
	setupOrcs(1)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
	checkContent(t, res.Body.String(), "Holding DoomHammer")
}

func TestEditOrc(t *testing.T) {
	setupOrcs(1)
	id := 1

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/orcs/edit/%d", id), nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
	checkContent(t, res.Body.String(), "Reforge an orc")
}

func TestDeleteOrc(t *testing.T) {
	setupOrcs(1)
	id := 1

	// Delete orc
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/orcs/delete/%d", id), nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusFound, res.Code)

	// Verify orc does not exist
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/orcs/view/%d", id), nil)
	res = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, res.Code)
}
