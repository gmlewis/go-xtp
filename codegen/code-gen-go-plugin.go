package codegen

// GeneratedFiles represents the files in the generated code.
type GeneratedFiles map[string]string

// genGoPluginPDK generates Plugin PDK code to process plugin calls in Go.
func (c *Client) genGoPluginPDK() (GeneratedFiles, error) {
	m := GeneratedFiles{
		"build.sh":               buildShScript,
		c.CustTypesFilename:      c.CustTypes,
		c.CustTypesTestsFilename: c.CustTypesTests,
	}

	return m, nil
}

var buildShScript = `#!/bin/bash -e
xtp plugin build
`
