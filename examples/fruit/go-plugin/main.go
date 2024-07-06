//go:build tinygo

// go-plugin represents an XTP Extension Plugin.
package main

import (
	"fmt"
	"time"

	"github.com/extism/go-pdk"
)

// VoidFunc - This demonstrates how you can create an export with
// no inputs or outputs.
func VoidFunc() {
	pdk.Log(pdk.LogDebug, "ENTER TinyGo plugin VoidFunc")
	// TODO: fill out your implementation here
	pdk.Log(pdk.LogDebug, "LEAVE TinyGo plugin VoidFunc")
}

// PrimitiveTypeFunc - This demonstrates how you can accept or return primtive types.
// This function takes a utf8 string and returns a json encoded boolean
//
// `input` - A string passed into plugin input
// Returns A boolean encoded as json
func PrimitiveTypeFunc(input string) bool {
	pdk.Log(pdk.LogDebug, fmt.Sprintf("ENTER TinyGo plugin PrimitiveTypeFunc('%v')", input))
	// TODO: fill out your implementation here
	pdk.Log(pdk.LogDebug, fmt.Sprintf("LEAVE TinyGo plugin PrimitiveTypeFunc('%v')", input))
	return false
}

// ReferenceTypeFunc - This demonstrates how you can accept or return references to schema types.
// And it shows how you can define an enum to be used as a property or input/output.
func ReferenceTypeFunc(input Fruit) ComplexObject {
	pdk.Log(pdk.LogDebug, fmt.Sprintf("ENTER TinyGo plugin ReferenceTypeFunc('%v')", input))
	// TODO: fill out your implementation here
	pdk.Log(pdk.LogDebug, fmt.Sprintf("LEAVE TinyGo plugin ReferenceTypeFunc('%v')", input))
	now := time.Now().Format(time.RFC3339)
	return ComplexObject{
		Ghost:          GhostGangEnumBlinky,
		ABoolean:       true,
		AString:        string(input),
		AnInt:          42,
		AnOptionalDate: &now,
	}
}

func main() {}
