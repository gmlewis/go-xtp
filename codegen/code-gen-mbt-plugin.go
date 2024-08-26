package codegen

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

var (
	mbtPluginHostFunctionsTemplate   = template.Must(template.New("cost-gen-mbt-plugin.go:mbtPluginHostFunctionsTemplateStr").Funcs(funcMap).Parse(mbtPluginHostFunctionsTemplateStr))
	mbtPluginMainTemplate            = template.Must(template.New("cost-gen-mbt-plugin.go:mbtPluginMainTemplateStr").Funcs(funcMap).Parse(mbtPluginMainTemplateStr))
	mbtPluginMoonPkgJSONTemplate     = template.Must(template.New("cost-gen-mbt-plugin.go:mbtPluginMoonPkgJSONTemplateStr").Funcs(funcMap).Parse(mbtPluginMoonPkgJSONTemplateStr))
	mbtPluginPluginFunctionsTemplate = template.Must(template.New("cost-gen-mbt-plugin.go:mbtPluginPluginFunctionsTemplateStr").Funcs(funcMap).Parse(mbtPluginPluginFunctionsTemplateStr))
	mbtPluginXtpTOMLTemplate         = template.Must(template.New("cost-gen-mbt-plugin.go:mbtPluginXtpTOMLTemplateStr").Parse(mbtPluginXtpTOMLTemplateStr))
)

// genMbtPluginPDK generates Plugin PDK code to process plugin calls in Mbt.
func (c *Client) genMbtPluginPDK() (GeneratedFiles, error) {
	var xtpTomlStr bytes.Buffer
	if err := mbtPluginXtpTOMLTemplate.Execute(&xtpTomlStr, c); err != nil {
		return nil, err
	}
	var hostFunctionsStr bytes.Buffer
	if err := mbtPluginHostFunctionsTemplate.Execute(&hostFunctionsStr, c); err != nil {
		return nil, err
	}
	var mainStr bytes.Buffer
	if err := mbtPluginMainTemplate.Execute(&mainStr, c); err != nil {
		return nil, err
	}
	var moonPkgJSONStr bytes.Buffer
	if err := mbtPluginMoonPkgJSONTemplate.Execute(&moonPkgJSONStr, c); err != nil {
		return nil, err
	}
	var pluginFunctionsStr bytes.Buffer
	if err := mbtPluginPluginFunctionsTemplate.Execute(&pluginFunctionsStr, c); err != nil {
		return nil, err
	}

	m := GeneratedFiles{
		"build.sh":             buildShScript,
		c.CustTypesFilename:    c.CustTypes,
		"main.mbt":             mainStr.String(),
		"moon.pkg.json":        moonPkgJSONStr.String(),
		"plugin-functions.mbt": pluginFunctionsStr.String(),
		"xtp.toml":             xtpTomlStr.String(),
	}

	if strings.TrimSpace(hostFunctionsStr.String()) != "" {
		m["host-functions.mbt"] = hostFunctionsStr.String()
	}

	return m, nil
}

//go:embed mbt-plugin-xtp-toml-template.txt
var mbtPluginXtpTOMLTemplateStr string

//go:embed mbt-plugin-host-functions-template.txt
var mbtPluginHostFunctionsTemplateStr string

//go:embed mbt-plugin-main-template.txt
var mbtPluginMainTemplateStr string

//go:embed mbt-plugin-moon-pkg-json-template.txt
var mbtPluginMoonPkgJSONTemplateStr string

// TODO: support primitives other than Strings.
//
//go:embed mbt-plugin-plugin-functions-template.txt
var mbtPluginPluginFunctionsTemplateStr string
