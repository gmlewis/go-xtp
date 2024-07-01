package codegen

import (
	"bytes"
	"text/template"
)

// GeneratedFiles represents the files in the generated code.
type GeneratedFiles map[string]string

var (
	goPluginXtpTOMLTemplate = template.Must(template.New("cost-gen-go-plugin.go:goPluginXtpTOMLTemplateStr").Parse(goPluginXtpTOMLTemplateStr))
)

// genGoPluginPDK generates Plugin PDK code to process plugin calls in Go.
func (c *Client) genGoPluginPDK() (GeneratedFiles, error) {
	var buf bytes.Buffer
	if err := goPluginXtpTOMLTemplate.Execute(&buf, c); err != nil {
		return nil, err
	}

	m := GeneratedFiles{
		"build.sh":               buildShScript,
		c.CustTypesFilename:      "package main\n\n" + c.CustTypes,
		c.CustTypesTestsFilename: "package main\n\n" + c.CustTypesTests,
		"xtp.toml":               buf.String(),
	}

	return m, nil
}

var buildShScript = `#!/bin/bash -e
xtp plugin build
`

var goPluginXtpTOMLTemplateStr = `app_id = "app_<enter-app-id-here>"

# This is where 'xtp plugin push' expects to find the wasm file after the build script has run.
bin = "{{ .PkgName }}.wasm"
extension_point_id = "ext_<enter-extension-point-id-here>"
name = "go-xtp-plugin-{{ .PkgName }}"

[scripts]

  # xtp plugin build runs this script to generate the wasm file
  build = "tinygo build -target wasi -o {{ .PkgName }}.wasm ."
`
