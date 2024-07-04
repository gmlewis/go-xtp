//go:build tinygo

package main

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

//export processUser
func processUser() int {
	input := pdk.InputString()
	v, err := ParseUser(input)
	if err != nil {
		pdk.Log(pdk.LogError, fmt.Sprintf("unable to ParseUser input: %v, input:\n%v\n", err, input))
		return 1 // failure
	}

	output := ProcessUser(v)

	buf, err := json.Marshal(output)
	if err != nil {
		pdk.Log(pdk.LogError, fmt.Sprintf("unable to json.Marshal output: %v", err))
		return 1 // failure
	}

	pdk.OutputString(string(buf))
	return 0 // success
}
