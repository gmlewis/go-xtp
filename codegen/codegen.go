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

// ClientOpts represents options to the codegen Client.
type ClientOpts struct {
	// Force causes existing file to be overwritten.
	Force bool
	// Quiet prevents warning messages from being printed
	Quiet bool
}

// Client represents a codegen client.
type Client struct {
	PkgName string
	Lang    string // "go" or "mbt"
	Plugin  *schema.Plugin

	CustTypesFilename string
	CustTypes         string

	CustTypesTestsFilename string
	CustTypesTests         string

	// internal fields used by the code generator:
	opts       ClientOpts
	numStructs int
}

// New returns a new codegen `Client` for either "go" or "mbt" and the
// provided plugin with the given package name.
func New(language string, plugin *schema.Plugin, opts *ClientOpts) (*Client, error) {
	if plugin == nil {
		return nil, errors.New("plugin cannot be nil")
	}
	if plugin.Version == "v0" {
		return nil, ErrNoCodeGeneration
	}
	if language != "go" && language != "mbt" {
		return nil, errors.New("language must be 'go' or 'mbt'")
	}
	if plugin.PkgName == "" {
		return nil, errors.New("plugin.PkgName must be provided")
	}

	c := &Client{
		PkgName: plugin.PkgName,
		Lang:    language,
		Plugin:  plugin,
	}
	if opts != nil {
		c.opts = *opts
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
