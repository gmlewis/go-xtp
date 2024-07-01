package codegen

import (
	"bytes"
	"html/template"
)

var (
	mbtPluginXtpTOMLTemplate = template.Must(template.New("cost-gen-mbt-plugin.go:mbtPluginXtpTOMLTemplateStr").Parse(mbtPluginXtpTOMLTemplateStr))
)

// genMbtPluginPDK generates Plugin PDK code to process plugin calls in Mbt.
func (c *Client) genMbtPluginPDK() (GeneratedFiles, error) {
	var buf bytes.Buffer
	if err := mbtPluginXtpTOMLTemplate.Execute(&buf, c); err != nil {
		return nil, err
	}

	m := GeneratedFiles{
		"build.sh":               buildShScript,
		c.CustTypesFilename:      c.CustTypes,
		c.CustTypesTestsFilename: c.CustTypesTests,
		"host-functions.mbt":     "", // TODO
		"main.mbt":               "", // TODO
		"moon.pkg.json":          "", // TODO
		"plugin-functions.mbt":   "", // TODO
		"xtp.toml":               buf.String(),
	}

	return m, nil
}

var mbtPluginXtpTOMLTemplateStr = `app_id = "app_<enter-app-id-here>"

# This is where 'xtp plugin push' expects to find the wasm file after the build script has run.
bin = "{{ .PkgName }}.wasm"
extension_point_id = "ext_<enter-extension-point-id-here>"
name = "mbt-xtp-plugin-{{ .PkgName }}"

[scripts]

  # xtp plugin build runs this script to generate the wasm file
  build = "moon build --target wasm && cp ../../../target/wasm/release/build/examples/{{ .PkgName }}/mbt-plugin/mbt-plugin.wasm ./{{ .PkgName }}.wasm"
`
