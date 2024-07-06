package codegen

import (
	"bytes"
	"text/template"
)

// GeneratedFiles represents the files in the generated code.
type GeneratedFiles map[string]string

var (
	goPluginHostFunctionsTemplate   = template.Must(template.New("cost-gen-go-plugin.go:goPluginHostFunctionsTemplateStr").Funcs(funcMap).Parse(goPluginHostFunctionsTemplateStr))
	goPluginMainTemplate            = template.Must(template.New("cost-gen-go-plugin.go:goPluginMainTemplateStr").Funcs(funcMap).Parse(goPluginMainTemplateStr))
	goPluginPluginFunctionsTemplate = template.Must(template.New("cost-gen-go-plugin.go:goPluginPluginFunctionsTemplateStr").Funcs(funcMap).Parse(goPluginPluginFunctionsTemplateStr))
	goPluginXtpTOMLTemplate         = template.Must(template.New("cost-gen-go-plugin.go:goPluginXtpTOMLTemplateStr").Parse(goPluginXtpTOMLTemplateStr))
)

// genGoPluginPDK generates Plugin PDK code to process plugin calls in Go.
func (c *Client) genGoPluginPDK() (GeneratedFiles, error) {
	var xtpTomlStr bytes.Buffer
	if err := goPluginXtpTOMLTemplate.Execute(&xtpTomlStr, c); err != nil {
		return nil, err
	}
	var hostFunctionsStr bytes.Buffer
	if err := goPluginHostFunctionsTemplate.Execute(&hostFunctionsStr, c); err != nil {
		return nil, err
	}
	var mainStr bytes.Buffer
	if err := goPluginMainTemplate.Execute(&mainStr, c); err != nil {
		return nil, err
	}
	var pluginFunctionsStr bytes.Buffer
	if err := goPluginPluginFunctionsTemplate.Execute(&pluginFunctionsStr, c); err != nil {
		return nil, err
	}

	m := GeneratedFiles{
		"build.sh":               buildShScript,
		c.CustTypesFilename:      "package main\n\n" + c.CustTypes,
		c.CustTypesTestsFilename: "package main\n\n" + c.CustTypesTests,
		"main.go":                mainStr.String(),
		"plugin-functions.go":    pluginFunctionsStr.String(),
		"xtp.toml":               xtpTomlStr.String(),
	}

	if len(c.Plugin.Imports) > 0 {
		m["host-functions.go"] = hostFunctionsStr.String()
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

var goPluginHostFunctionsTemplateStr = `//go:build tinygo

package main

import (
	"encoding/json"

	"github.com/extism/go-pdk"
)
{{range .Plugin.Imports }}{{ $name := .Name }}
//go:wasmimport extism:host/user {{ $name }}
func host{{ $name | uppercaseFirst }}(uint64) uint64

// {{ $name | uppercaseFirst }} - {{ .Description | goMultilineComment | stripLeadingSlashes | leftJustify }}
func {{ $name | uppercaseFirst }}({{ .Input | inputToGoType }}) ({{ .Output | outputToGoType }}, error) {
	buf, err := json.Marshal(input)
	if err != nil {
		return false, err
	}

	mem := pdk.AllocateBytes(buf)
	ptr := host{{ $name | uppercaseFirst }}(mem.Offset())

	rmem := pdk.FindMemory(ptr)
	buf = rmem.ReadBytes()

	var result {{ .Output | jsonOutputAsGoType }}
	if err := json.Unmarshal(buf, &result); err != nil {
		return false, err
	}
	return result, nil
}
{{ end }}`

var goPluginMainTemplateStr = `//go:build tinygo

// go-plugin represents an XTP Extension Plugin.
package main

import "github.com/extism/go-pdk"
{{range .Plugin.Exports }}{{ $name := .Name }}
// {{ $name | uppercaseFirst }} - {{ .Description | goMultilineComment | stripLeadingSlashes | leftJustify }}{{ if exportHasInputOrOutputDescription . }}
//
{{ end }}{{ if exportHasInputDescription . }}// ` + "`input`" + ` - {{ .Input.Description | goMultilineComment | stripLeadingSlashes | leftJustify }}{{ end }}{{ if exportHasOutputDescription . }}
// Returns {{ .Output.Description | goMultilineComment | stripLeadingSlashes | leftJustify }}{{ end }}
func {{ $name | uppercaseFirst }}({{ .Input | inputToGoType }}){{ if .Output }} {{ .Output | outputToGoType }}{{ end }} {
{{ "\t" }}pdk.Log(pdk.LogDebug, "ENTER TinyGo plugin {{ $name | uppercaseFirst }}")
{{ "\t" }}// TODO: fill out your implementation here
{{ "\t" }}pdk.Log(pdk.LogDebug, "LEAVE TinyGo plugin {{ $name | uppercaseFirst }}"){{ .Output | outputToGoExampleLiteral }}
}
{{ end }}
func main() {}
`

// TODO: support primitive types other than string.
var goPluginPluginFunctionsTemplateStr = `//go:build tinygo

package main

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)
{{range $index, $export := .Plugin.Exports }}{{ $name := .Name }}
//export {{ $name }}
func {{ $name }}() int {
{{ if . | inputIsVoidType }}	{{ $name | uppercaseFirst }}(){{ end -}}
{{ if . | inputIsPrimitiveType }}	var input string
	if err := json.Unmarshal([]byte(pdk.InputString()), &input); err != nil {
		pdk.Log(pdk.LogError, fmt.Sprintf("unable to json.Unmarshal input: %v", err))
		return 1 // failure
	}

	output := {{ $name | uppercaseFirst }}(input)

	buf, err := json.Marshal(output)
	if err != nil {
		pdk.Log(pdk.LogError, fmt.Sprintf("unable to json.Marshal output: %v", err))
		return 1 // failure
	}

	pdk.OutputString(string(buf)){{ end -}}
{{ if . | inputIsReferenceType }}	input := pdk.InputString()
	v, err := Parse{{ inputReferenceTypeName . }}(input)
	if err != nil {
		pdk.Log(pdk.LogError, fmt.Sprintf("unable to Parse{{ inputReferenceTypeName . }} input: %v, input:\n%v\n", err, input))
		return 1 // failure
	}

	output := {{ $name | uppercaseFirst }}(v)

	buf, err := json.Marshal(output)
	if err != nil {
		pdk.Log(pdk.LogError, fmt.Sprintf("unable to json.Marshal output: %v", err))
		return 1 // failure
	}

	pdk.OutputString(string(buf)){{ end }}
	return 0 // success
{{ "}" }}
{{ end }}`
