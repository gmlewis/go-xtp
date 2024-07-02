package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

// Fruit represents a set of available fruits you can consume.
type Fruit string

const (
	FruitEnumApple      Fruit = "apple"
	FruitEnumOrange     Fruit = "orange"
	FruitEnumBanana     Fruit = "banana"
	FruitEnumStrawberry Fruit = "strawberry"
)

// ParseFruit parses a JSON string and returns the value.
func ParseFruit(s string) (value Fruit, err error) {
	switch s {
	case "apple":
		return FruitEnumApple, nil
	case "orange":
		return FruitEnumOrange, nil
	case "banana":
		return FruitEnumBanana, nil
	case "strawberry":
		return FruitEnumStrawberry, nil
	default:
		return value, fmt.Errorf("not a Fruit: %v", s)
	}
}

// GhostGang represents a set of all the enemies of pac-man.
type GhostGang string

const (
	GhostGangEnumBlinky GhostGang = "blinky"
	GhostGangEnumPinky  GhostGang = "pinky"
	GhostGangEnumInky   GhostGang = "inky"
	GhostGangEnumClyde  GhostGang = "clyde"
)

// ParseGhostGang parses a JSON string and returns the value.
func ParseGhostGang(s string) (value GhostGang, err error) {
	switch s {
	case "blinky":
		return GhostGangEnumBlinky, nil
	case "pinky":
		return GhostGangEnumPinky, nil
	case "inky":
		return GhostGangEnumInky, nil
	case "clyde":
		return GhostGangEnumClyde, nil
	default:
		return value, fmt.Errorf("not a GhostGang: %v", s)
	}
}

// ComplexObject represents a complex json object.
type ComplexObject struct {
	// I can override the description for the property here
	Ghost GhostGang `json:"ghost"`
	// A boolean prop
	ABoolean bool `json:"aBoolean"`
	// An string prop
	AString string `json:"aString"`
	// An int prop
	AnInt int `json:"anInt"`
	// A datetime object, we will automatically serialize and deserialize
	// this for you.
	AnOptionalDate *string `json:"anOptionalDate,omitempty"`
}

// ParseComplexObject parses a JSON string and returns the value.
func ParseComplexObject(s string) (value ComplexObject, err error) {
	if err := jsoncomp.Unmarshal([]byte(s), &value); err != nil {
		return value, err
	}

	return value, nil
}

// GetSchema returns an `XTPSchema` for the `ComplexObject`.
func (c *ComplexObject) GetSchema() XTPSchema {
	return XTPSchema{
		"ghost":          "GhostGang",
		"aBoolean":       "boolean",
		"aString":        "string",
		"anInt":          "integer",
		"anOptionalDate": "?Date",
	}
}

// XTPSchema describes the values and types of an XTP object
// in a language-agnostic format.
type XTPSchema map[string]string
