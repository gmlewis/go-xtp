//go:build tinygo

package main

import (
	"encoding/json"
	"fmt"

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
		pdk.Log(pdk.LogError, fmt.Errorf("unable to json.Marshal output: %v", err))
		return 1 // failure
	}

	pdk.OutputString(string(buf))
	return 0 // success
}

//export referenceTypeFunc
func referenceTypeFunc() int {
	input := pdk.InputString()
	v, err := ParseFruit(input)
	if err != nil {
		pdk.Log(pdk.LogError, fmt.Errorf("unable to ParseFruit input: %v, input:\n%v\n", err, input))
		return 1 // failure
	}

	output := ReferenceTypeFunc(v)

	buf, err := json.Marshal(output)
	if err != nil {
		pdk.Log(pdk.LogError, fmt.Errorf("unable to json.Marshal output: %v", err))
		return 1 // failure
	}

	pdk.OutputString(string(buf))
	return 0 // success
}
