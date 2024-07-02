//go:build tinygo

package main

import (
	"encoding/json"

	"github.com/extism/go-pdk"
)

//export processUser
func processUser() int {
	input := pdk.InputString()
	output := ProcessUser(User(input))

	buf, err := json.Marshal(output)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return 1 // failure
	}

	pdk.OutputString(string(buf))
	return 0 // success
}
