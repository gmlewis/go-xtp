// Package user represents the custom datatypes for an XTP Extension Plugin.
package user

// Address represents a users address.
type Address struct {
	// Street address
	Street string `json:"street"`
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

// GetSchema returns an `XTPSchema` for the `User`.
func (c *User) GetSchema() XTPSchema {
	return XTPSchema{
		"age":     "?integer",
		"email":   "?string",
		"address": "?Address",
	}
}

// XTPSchema describes the values and types of an XTP object
// in a language-agnostic format.
type XTPSchema map[string]string
