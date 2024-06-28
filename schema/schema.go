// Package schema provides methods to parse the XTP Extension Plugin schema.
package schema

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	versionRE = regexp.MustCompile(`(?m)^version: v(\d+)`)
)

// Plugin represents an XTP Extension Plugin Schema.
type Plugin struct {
	Version     string        `yaml:"version"`
	Exports     []*Export     `yaml:"exports"`
	Imports     []*Import     `yaml:"imports,omitempty"`
	CustomTypes []*CustomType `yaml:"schemas,omitempty"`

	// the following fields are only used by the code generator:
	Lang string `yaml:"-"`
}

// Export represents an exported function by the XTP Extension Plugin.
type Export struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description,omitempty"`
	Input       *Input        `yaml:"input,omitempty"`
	Output      *Output       `yaml:"output,omitempty"`
	CodeSamples []*CodeSample `yaml:"codeSamples,omitempty"`
}

// Input represents an input to the exported function.
type Input struct {
	Ref         string `yaml:"$ref,omitempty"`
	Type        string `yaml:"type,omitempty"`
	Description string `yaml:"description,omitempty"`
	ContentType string `yaml:"contentType,omitempty"`
}

// Output represents an output from the exported function.
type Output struct {
	Ref         string `yaml:"$ref,omitempty"`
	Type        string `yaml:"type,omitempty"`
	Description string `yaml:"description,omitempty"`
	ContentType string `yaml:"contentType,omitempty"`
}

// CodeSample represents a code sample for calling the function in a
// particular programming language.
type CodeSample struct {
	Lang   string `yaml:"lang"`
	Label  string `yaml:"label,omitempty"`
	Source string `yaml:"source"`
}

// Import represents an imported function into the XTP Extension Plugin from the host.
type Import struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description,omitempty"`
	Input       *Input  `yaml:"input,omitempty"`
	Output      *Output `yaml:"output,omitempty"`
}

// CustomType represents an XTP Extension Plugin custom datatype.
type CustomType struct {
	Name        string      `yaml:"name"`
	ContentType string      `yaml:"contentType,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Enum        []string    `yaml:"enum,omitempty"`
	Required    []string    `yaml:"required,omitempty"`
	Properties  []*Property `yaml:"properties,omitempty"`
}

// GetRequiredProps returns the required properties for this CustomType.
func (ct *CustomType) GetRequiredProps() []*Property {
	reqFields := map[string]bool{}
	for _, req := range ct.Required {
		reqFields[req] = true
	}
	reqProps := make([]*Property, 0, len(reqFields))
	for _, prop := range ct.Properties {
		if reqFields[prop.Name] {
			reqProps = append(reqProps, prop)
		}
	}
	return reqProps
}

// Property represents an argument to a plugin function.
type Property struct {
	Name        string   `yaml:"name"`
	Ref         string   `yaml:"$ref,omitempty"`
	Type        string   `yaml:"type,omitempty"`
	Format      string   `yaml:"format,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Maximum     *float64 `yaml:"maximum,omitempty"`
	Minimum     *float64 `yaml:"minimum,omitempty"`

	// the following fields are only used by the code generator:
	FirstEnumValue string      `yaml:"-"`
	IsRequired     bool        `yaml:"-"`
	RefCustomType  *CustomType `yaml:"-"`
}

// ParseStr parses an XTP Extension Plugin schema yaml string and returns it.
func ParseStr(yamlStr string) (*Plugin, error) {
	m := versionRE.FindStringSubmatch(yamlStr)
	if len(m) != 2 {
		return nil, errors.New("unable to find schema version")
	}

	switch m[1] {
	case "0":
		return ParseV0(yamlStr)
	case "1":
		return ParseV1(yamlStr)
	default:
		return nil, fmt.Errorf("unsupported yaml version %v", m[1])
	}
}
