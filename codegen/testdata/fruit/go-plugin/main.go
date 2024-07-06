//go:build tinygo

// go-plugin represents an XTP Extension Plugin.
package main

import "github.com/extism/go-pdk"

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
	pdk.Log(pdk.LogDebug, "ENTER TinyGo plugin PrimitiveTypeFunc")
	// TODO: fill out your implementation here
	pdk.Log(pdk.LogDebug, "LEAVE TinyGo plugin PrimitiveTypeFunc")
	return false
}

// ReferenceTypeFunc - This demonstrates how you can accept or return references to schema types.
// And it shows how you can define an enum to be used as a property or input/output.
func ReferenceTypeFunc(input Fruit) ComplexObject {
	pdk.Log(pdk.LogDebug, "ENTER TinyGo plugin ReferenceTypeFunc")
	// TODO: fill out your implementation here
	pdk.Log(pdk.LogDebug, "LEAVE TinyGo plugin ReferenceTypeFunc")
	return ComplexObject{}
}

func main() {}
