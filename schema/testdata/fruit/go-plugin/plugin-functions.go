package main

import (
	"encoding/json"

	"github.com/extism/go-pdk"
)

//export voidFunc
func voidFunc() int {
	VoidFunc()
	return 0 // success
}

//export primitiveTypeFunc
func primitiveTypeFunc() int {
	input := pdk.InputString()
	output := PrimitiveTypeFunc(input)

	buf, err := json.Marshal(output)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return 1 // failure
	}

	pdk.OutputString(string(buf))
	return 0 // success
}

//export referenceTypeFunc
func referenceTypeFunc() int {
	input := pdk.InputString()
	output := ReferenceTypeFunc(Fruit(input))

	buf, err := json.Marshal(output)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return 1 // failure
	}

	pdk.OutputString(string(buf))
	return 0 // success
}
