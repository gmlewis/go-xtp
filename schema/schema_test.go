package schema

import (
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/fruit.yaml
var fruitYaml string

//go:embed testdata/user.yaml
var userYaml string

//go:embed testdata/v0.yaml
var v0Yaml string

func floatPtr(f float64) *float64 { return &f }

func TestParseStr(t *testing.T) {
	t.Parallel()

	// To test the "user" struct, we need to point two places in the test to the same struct, so define it here:
	addressCustomTypePointer := &CustomType{
		Name:        "Address",
		Description: "A users address",
		Required:    []string{"street"},
		Properties: []*Property{
			{Name: "street", Description: "Street address", Type: "string", IsRequired: true},
		},
		ContentType: "application/json",
	}

	tests := []struct {
		name    string
		yamlStr string
		want    *Plugin
	}{
		{
			name:    "fruit",
			yamlStr: fruitYaml,
			want: &Plugin{
				Version: "v1-draft",
				Exports: []*Export{
					{
						Name:        "voidFunc",
						Description: "This demonstrates how you can create an export with\nno inputs or outputs.\n",
					},
					{
						Name:        "primitiveTypeFunc",
						Description: "This demonstrates how you can accept or return primtive types.\nThis function takes a utf8 string and returns a json encoded boolean\n",
						Input: &Input{
							Type:        "string",
							Description: "A string passed into plugin input",
							ContentType: "text/plain; charset=UTF-8",
						},
						Output: &Output{
							Type:        "boolean",
							Description: "A boolean encoded as json",
							ContentType: "application/json",
						},
						CodeSamples: []*CodeSample{
							{
								Lang:  "typescript",
								Label: "Test if a string has more than one character.\nCode samples show up in documentation and inline in docstrings\n",
								Source: `function primitiveTypeFunc(input: string): boolean {
  return input.length > 1
}
`,
							},
						},
					},
					{
						Name:        "referenceTypeFunc",
						Description: "This demonstrates how you can accept or return references to schema types.\nAnd it shows how you can define an enum to be used as a property or input/output.\n",
						Input: &Input{
							Ref: "#/schemas/Fruit",
						},
						Output: &Output{
							Ref: "#/schemas/ComplexObject",
						},
					},
				},
				Imports: []*Import{
					{
						Name:        "eatAFruit",
						Description: "This is a host function. Right now host functions can only be the type (i64) -> i64.\nWe will support more in the future. Much of the same rules as exports apply.\n",
						Input:       &Input{Ref: "#/schemas/Fruit"},
						Output: &Output{
							Type:        "boolean",
							Description: "boolean encoded as json",
							ContentType: "application/json",
						},
					},
				},
				CustomTypes: []*CustomType{
					{
						Name:        "Fruit",
						Description: "A set of available fruits you can consume",
						Enum:        []string{"apple", "orange", "banana", "strawberry"},
					},
					{
						Name:        "GhostGang",
						Description: "A set of all the enemies of pac-man",
						Enum:        []string{"blinky", "pinky", "inky", "clyde"},
					},
					{
						Name:        "ComplexObject",
						Description: "A complex json object",
						Required:    []string{"ghost", "aBoolean", "aString", "anInt"},
						Properties: []*Property{
							{
								Name:           "ghost",
								Ref:            "#/schemas/GhostGang",
								Description:    "I can override the description for the property here",
								IsRequired:     true,
								FirstEnumValue: "blinky",
							},
							{Name: "aBoolean", Description: "A boolean prop", Type: "boolean", IsRequired: true},
							{Name: "aString", Description: "An string prop", Type: "string", IsRequired: true},
							{Name: "anInt", Description: "An int prop", Type: "integer", Format: "int32", IsRequired: true},
							{
								Name:        "anOptionalDate",
								Description: "A datetime object, we will automatically serialize and deserialize\nthis for you.\n",
								Type:        "string",
								Format:      "date-time",
							},
						},
						ContentType: "application/json",
					},
				},
			},
		},
		{
			name:    "user",
			yamlStr: userYaml,
			want: &Plugin{
				Version: "v1-draft",
				Exports: []*Export{
					{
						Name:        "processUser",
						Description: "The second export function",
						Input:       &Input{Ref: "#/schemas/User"},
						Output:      &Output{Ref: "#/schemas/User"},
						CodeSamples: []*CodeSample{
							{
								Lang:  "typescript",
								Label: "Process a user by email",
								Source: `function processUser(user: User): User {
  if (user.email.endsWith('@aol.com')) user.age += 10
  return user
}
`,
							},
						},
					},
				},
				CustomTypes: []*CustomType{
					addressCustomTypePointer,
					{
						Name:        "User",
						Description: "A user object in our system.",
						Properties: []*Property{
							{
								Name:        "age",
								Description: "The user's age, naturally",
								Type:        "integer",
								Format:      "int32",
								Maximum:     floatPtr(200),
								Minimum:     floatPtr(0),
							},
							{
								Name:        "email",
								Description: "The user's email, of course", Type: "string",
							},
							{
								Name:          "address",
								Ref:           "#/schemas/Address",
								RefCustomType: addressCustomTypePointer,
							},
						},
						ContentType: "application/json",
					},
				},
			},
		},
		{
			name:    "v0",
			yamlStr: v0Yaml,
			want: &Plugin{
				Version: "v0",
				Exports: []*Export{
					{Name: "myExport"},
					{Name: "processUser"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStr(tt.yamlStr)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ParseStr mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestToYaml(t *testing.T) {
	tests := []struct {
		name    string
		yamlStr string
	}{
		{
			name:    "fruit",
			yamlStr: fruitYaml,
		},
		{
			name:    "user",
			yamlStr: userYaml,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin, err := ParseStr(tt.yamlStr)
			if err != nil {
				t.Fatal(err)
			}

			got, err := plugin.ToYaml()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.yamlStr, got); diff != "" {
				t.Errorf("ParseStr mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
