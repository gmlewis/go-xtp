// Package user represents the custom datatypes for an XTP Extension Plugin.
package user

// Address represents a users address.
type Address struct {
	// Street address
	Street string `json:"street"`
}

// User represents a user object in our system..
type User struct {
	// The user's age, naturally
	Age *int `json:"age,omitempty"`
	// The user's email, of course
	Email   *string  `json:"email,omitempty"`
	Address *Address `json:"address,omitempty"`
}

