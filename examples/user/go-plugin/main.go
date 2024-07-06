//go:build tinygo

// go-plugin represents an XTP Extension Plugin.
package main

import (
	"fmt"

	"github.com/extism/go-pdk"
)

// ProcessUser - The second export function
func ProcessUser(input User) User {
	pdk.Log(pdk.LogDebug, fmt.Sprintf("ENTER TinyGo plugin ProcessUser('%v')", input))
	// TODO: fill out your implementation here
	pdk.Log(pdk.LogDebug, fmt.Sprintf("LEAVE TinyGo plugin ProcessUser('%v')", input))
	age := 42
	email := "email@example.com"
	return User{
		Age:     &age,
		Email:   &email,
		Address: &Address{Street: "main street"},
	}
}

func main() {}
