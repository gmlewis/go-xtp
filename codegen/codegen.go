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
	PkgName string
	Lang    string // "go" or "mbt"
	Plugin  *schema.Plugin

	CustTypesFilename string
	CustTypes         string

	CustTypesTestsFilename string
	CustTypesTests         string

	force bool
}

// New returns a new codegen `Client` for either "go" or "mbt" and the
// provided plugin with the given package name.
func New(language, packageName string, plugin *schema.Plugin, force bool) (*Client, error) {
	if plugin == nil {
		return nil, errors.New("plugin cannot be nil")
	}
	if plugin.Version == "v0" {
		return nil, ErrNoCodeGeneration
	}
	if language != "go" && language != "mbt" {
		return nil, errors.New("language must be 'go' or 'mbt'")
	}
	if packageName == "" {
		return nil, errors.New("packageName must be provided")
	}

	c := &Client{
		PkgName: packageName,
		Lang:    language,
		Plugin:  plugin,
		force:   force,
	}

	if err := c.genCustomTypes(); err != nil {
		return nil, err
	}

	if c.CustTypes == "" || c.CustTypesFilename == "" {
		return nil, errors.New("programming error: CustTypes or CustTypesFilename empty")
	}
	if c.CustTypesTests == "" || c.CustTypesTestsFilename == "" {
		return nil, errors.New("programming error: CustTypesTests or CustTypesTestsFilename empty")
	}

	return c, nil
}
