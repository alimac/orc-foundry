package orc

import (
	"bytes"
	"strings"
	"testing"
)

// generateCombinations runs in O(n^2) time
// TODO figure out a more efficient way
func generateCombinations(a []string, b []string, result string) bool {
	// Use a buffer instead of concatenating string with +
	// http://programming-tips.info/concatenate_string_literal/go/en/index.html
	// https://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go
	var buffer bytes.Buffer

	for _, i := range a {
		for _, j := range b {
			buffer.WriteString(i)
			buffer.WriteString(j)
			if strings.ToLower(result) == buffer.String() {
				return true
			}
			buffer.Reset()
		}
	}
	return false
}

func TestOrcGreeting(t *testing.T) {
	greeting := Forge("greeting")
	var found = false

	for _, item := range greetings {
		if item == strings.ToLower(greeting) {
			found = true
			t.Logf("Greeting = %s is a valid greeting", greeting)
		}
	}

	if !found {
		t.Errorf("Greeting = %s is not found", greeting)
	}
}

func TestOrcName(t *testing.T) {
	name := Forge("name")
	found := generateCombinations(prefixes, suffixes, name)

	if !found {
		t.Errorf("Name = %s is not a valid combination", name)
	} else {
		t.Logf("Name = %s is a valid name", name)
	}
}

func TestOrcWeapon(t *testing.T) {
	weapon := Forge("weapon")
	found := generateCombinations(adjectives, weapons, weapon)

	if !found {
		t.Errorf("Weapon = %s is not a valid combination", weapon)
	} else {
		t.Logf("Weapon = %s is a valid greeting", weapon)
	}
}

func TestDefault(t *testing.T) {
	result := Forge("")
	expected := ""
	if result != expected {
		t.Errorf("Default returned %s, want \"\"", result)
	}
}

func BenchmarkOrcName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Forge("name")
	}
}

func BenchmarkOrcGreeting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Forge("greeting")
	}
}

func BenchmarkOrcWeapon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Forge("weapon")
	}
}
