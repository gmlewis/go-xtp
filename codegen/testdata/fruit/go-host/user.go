package main

import (
	"encoding/json" // jsoniter/jsoncomp are not compatible with tinygo.
)

// Address represents a users address.
type Address struct {
	// Street address
	Street string `json:"street"`
}

// ParseAddress parses a JSON string and returns the value.
func ParseAddress(s string) (value Address, err error) {
	if err := json.Unmarshal([]byte(s), &value); err != nil {
		return value, err
	}

	return value, nil
}

// GetSchema returns an `XTPSchema` for the `Address`.
func (c *Address) GetSchema() XTPSchema {
	return XTPSchema{
		"street": "string",
	}
}

// User represents a user object in our system..
type User struct {
	// The user's age, naturally
	Age *int `json:"age,omitempty"`
	// The user's email, of course
	Email   *string  `json:"email,omitempty"`
	Address *Address `json:"address,omitempty"`
}

// ParseUser parses a JSON string and returns the value.
func ParseUser(s string) (value User, err error) {
	if err := json.Unmarshal([]byte(s), &value); err != nil {
		return value, err
	}

	return value, nil
}

// GetSchema returns an `XTPSchema` for the `User`.
func (c *User) GetSchema() XTPSchema {
	return XTPSchema{
		"age":     "?integer",
		"email":   "?string",
		"address": "?Address",
	}
}
