// go-plugin represents an XTP Extension Plugin.
package main

// VoidFunc - This demonstrates how you can create an export with
// no inputs or outputs.
func VoidFunc() {
	// TODO: fill out your implementation here
	panic("VoidFunc not yet implemented.")
}

// PrimitiveTypeFunc - This demonstrates how you can accept or return primtive types.
// This function takes a utf8 string and returns a json encoded boolean
//
// `input` - A string passed into plugin input
// Returns A boolean encoded as json
func PrimitiveTypeFunc(input string) bool {
	// TODO: fill out your implementation here
	panic("PrimitiveTypeFunc not yet implemented.")
}

// ReferenceTypeFunc - This demonstrates how you can accept or return references to schema types.
// And it shows how you can define an enum to be used as a property or input/output.
func ReferenceTypeFunc(input Fruit) ComplexObject {
	// TODO: fill out your implementation here
	panic("ReferenceTypeFunc not yet implemented.")
}

func main() {}
