package main

import (
	"encoding/json"

	"github.com/extism/go-pdk"
)

//go:wasmimport extism:host/user eatAFruit
func hostEatAFruit(uint64) uint64

// EatAFruit - This is a host function. Right now host functions can only be the type (i64) -> i64.
// We will support more in the future. Much of the same rules as exports apply.
func EatAFruit(input Fruit) (bool, error) {
	buf, err := json.Marshal(input)
	if err != nil {
		return false, err
	}

	mem := pdk.AllocateBytes(buf)
	ptr := hostEatAFruit(mem.Offset())

	rmem := pdk.FindMemory(ptr)
	buf = rmem.ReadBytes()

	var result bool
	if err := json.Unmarshal(buf, &result); err != nil {
		return false, err
	}
	return result, nil
}
