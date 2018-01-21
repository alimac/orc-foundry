// Package orc provides helper methods to generate orc attributes.
//
package orc

import (
	"math/rand"
	"strings"
	"time"
)

// Forge returns the requested orc attribute.
// Supported attributes: name, greeting, weapon.
func Forge(attribute string) string {
	rand.Seed(time.Now().UnixNano())

	switch attribute {
	case "name":
		prefix := strings.Title(prefixes[rand.Intn(len(prefixes))])
		suffix := suffixes[rand.Intn(len(suffixes))]
		return prefix + suffix
	case "greeting":
		return strings.Title(greetings[rand.Intn(len(greetings))])
	case "weapon":
		adjective := strings.Title(adjectives[rand.Intn(len(adjectives))])
		weapon := strings.Title(weapons[rand.Intn(len(weapons))])
		return adjective + weapon
	default:
		return ""
	}
}
