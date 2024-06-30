// Package codegen generates custom datatypes, PDK plugin code and SDK host code
// from a `schema.Plugin` in the Go and MoonBit programming languages.
package codegen

import (
	"errors"

	"github.com/gmlewis/go-xtp/schema"
)

var (
	ErrNoCodeGeneration = errors.New("code generation not supported for version v0")
)

// Client represents a codegen client.
type Client struct {
	Lang   string // "go" or "mbt"
	Plugin *schema.Plugin

	CustTypes      string
	CustTypesTests string
}

// New returns a new codegen `Client` for either "go" or "mbt" and the
// provided plugin.
func New(language string, plugin *schema.Plugin) (*Client, error) {
	if plugin == nil {
		return nil, errors.New("plugin cannot be nil")
	}
	if plugin.Version == "v0" {
		return nil, ErrNoCodeGeneration
	}
	if language != "go" && language != "mbt" {
		return nil, errors.New("language must be 'go' or 'mbt'")
	}

	c := &Client{
		Lang:   language,
		Plugin: plugin,
	}

	custTypes, custTypesTests, err := c.genCustomTypes()
	if err != nil {
		return nil, err
	}

	c.CustTypes = custTypes
	c.CustTypesTests = custTypesTests

	return c, nil
}
