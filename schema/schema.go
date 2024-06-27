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
	Version string    `yaml:"version"`
	Exports []*Export `yaml:"exports"`
	Schemas []*Func   `yaml:"schemas,omitempty"`
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
	Description string `yaml:"description,omitempty"`
	Type        string `yaml:"type,omitempty"`
	ContentType string `yaml:"contentType,omitempty"`
}

// Output represents an output from the exported function.
type Output struct {
	Ref         string `yaml:"$ref,omitempty"`
	Description string `yaml:"description,omitempty"`
	Type        string `yaml:"type,omitempty"`
	ContentType string `yaml:"contentType,omitempty"`
}

// CodeSample represents a code sample for calling the function in a
// particular programming language.
type CodeSample struct {
	Lang   string `yaml:"lang"`
	Label  string `yaml:"label,omitempty"`
	Source string `yaml:"source"`
}

// Func represents an XTP Extension Plugin function.
type Func struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description,omitempty"`
	Properties  []*Property `yaml:"properties,omitempty"`
	ContentType string      `yaml:"contentType,omitempty"`
}

// Property represents an argument to a plugin function.
type Property struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Type        string   `yaml:"type,omitempty"`
	Format      string   `yaml:"format,omitempty"`
	Maximum     *float64 `yaml:"maximum,omitempty"`
	Minimum     *float64 `yaml:"minimum,omitempty"`
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
