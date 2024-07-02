//go:build tinygo

// go-plugin represents an XTP Extension Plugin.
package main

// ProcessUser - The second export function
func ProcessUser(input User) User {
	return input
}

func main() {}
