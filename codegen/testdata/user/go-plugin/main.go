//go:build tinygo

// go-plugin represents an XTP Extension Plugin.
package main

import "github.com/extism/go-pdk"

// ProcessUser - The second export function
func ProcessUser(input User) User {
	pdk.Log(pdk.LogDebug, "ENTER TinyGo plugin ProcessUser")
	// TODO: fill out your implementation here
	pdk.Log(pdk.LogDebug, "LEAVE TinyGo plugin ProcessUser")
	return User{}
}

func main() {}
