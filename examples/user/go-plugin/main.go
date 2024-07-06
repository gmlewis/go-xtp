//go:build tinygo

// go-plugin represents an XTP Extension Plugin.
package main

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

// ProcessUser - The second export function
func ProcessUser(input User) User {
	inBuf, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		pdk.Log(pdk.LogError, fmt.Sprintf("MarshalIndent error: %v", err))
	}
	pdk.Log(pdk.LogDebug, fmt.Sprintf("ENTER TinyGo plugin ProcessUser(): input=\n%s", inBuf))
	// TODO: fill out your implementation here
	pdk.Log(pdk.LogDebug, "LEAVE TinyGo plugin ProcessUser")
	age := 42
	email := "email@example.com"
	return User{
		Age:     &age,
		Email:   &email,
		Address: &Address{Street: "main street"},
	}
}

func main() {}
