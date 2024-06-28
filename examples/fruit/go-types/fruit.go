// Package fruit represents the custom datatypes for an XTP Extension Plugin.
package fruit

// Fruit represents a set of available fruits you can consume.
type Fruit string

const (
	FruitEnumApple      Fruit = "apple"
	FruitEnumOrange     Fruit = "orange"
	FruitEnumBanana     Fruit = "banana"
	FruitEnumStrawberry Fruit = "strawberry"
)

// GhostGang represents a set of all the enemies of pac-man.
type GhostGang string

const (
	GhostGangEnumBlinky GhostGang = "blinky"
	GhostGangEnumPinky  GhostGang = "pinky"
	GhostGangEnumInky   GhostGang = "inky"
	GhostGangEnumClyde  GhostGang = "clyde"
)

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
