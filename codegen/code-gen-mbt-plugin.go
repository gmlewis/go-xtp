package codegen

import (
	"bytes"
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
		"build.sh":               buildShScript,
		c.CustTypesFilename:      c.CustTypes,
		c.CustTypesTestsFilename: c.CustTypesTests,
		"host-functions.mbt":     hostFunctionsStr.String(),
		"main.mbt":               mainStr.String(),
		"moon.pkg.json":          moonPkgJSONStr.String(),
		"plugin-functions.mbt":   pluginFunctionsStr.String(),
		"xtp.toml":               xtpTomlStr.String(),
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

var mbtPluginHostFunctionsTemplateStr = `{{range .Plugin.Imports }}{{ $name := .Name }}pub fn host_{{ $name | lowerSnakeCase }}(offset : Int64) -> Int64 = "extism:host/user" "{{ $name }}"

/// ` + "`{{ $name | lowerSnakeCase }}`" + ` - {{ .Description | mbtMultilineComment | stripLeadingSlashes | leftJustify }}
pub fn {{ $name | lowerSnakeCase }}({{ .Input | inputToMbtType }}) -> {{ .Output | outputToMbtType }}!String {
  let json = @jsonutil.to_json(input)
  let mem = @host.Memory::allocate_json_value(json)
  let ptr = host_{{ $name | lowerSnakeCase }}(mem.offset)
  let buf = @host.find_memory(ptr).to_string()
  let out = @json.parse(buf)
  match out {
    Ok(jv) =>
      match jv.{{ .Output | jsonOutputAsMbtType }}() {
        Some(v) => v
        None => {
          raise "unable to parse \(buf)"
        }
      }
    Err(e) => {
      raise "unable to parse \(buf): \(e)"
    }
  }
}{{ end }}
`

var mbtPluginMainTemplateStr = `{{range .Plugin.Exports }}{{ $name := .Name }}/// ` + "`{{ $name | lowerSnakeCase }}`" + ` - {{ .Description | mbtMultilineComment | stripLeadingSlashes | leftJustify }}{{ if exportHasInputOrOutputDescription . }}
///
{{ end }}{{ if exportHasInputDescription . }}/// ` + "`input`" + ` - {{ .Input.Description | mbtMultilineComment | stripLeadingSlashes | leftJustify }}{{ end }}{{ if exportHasOutputDescription . }}
/// Returns {{ .Output.Description | mbtMultilineComment | stripLeadingSlashes | leftJustify }}{{ end }}
pub fn {{ $name | lowerSnakeCase }}({{ .Input | inputToMbtType }}) -> {{ .Output | outputToMbtType }} {
  // TODO: fill out your implementation here{{ .Output | outputToMbtExampleLiteral }}
}

{{ end }}fn main {

}
`

var mbtPluginMoonPkgJSONTemplateStr = `{
  "is-main": true,
  "import": [
    "gmlewis/moonbit-pdk/pdk/host",
    {
      "path": "gmlewis/json",
      "alias": "jsonutil"
    }
  ],
  "link": {
    "wasm": {
      "exports": [{{ $exportsLen := .Plugin.Exports | len }}{{range $index, $export := .Plugin.Exports }}{{ $name := .Name }}
        "exported_{{ $name | lowerSnakeCase }}:{{ $name }}"{{ showJSONCommaForOptional $index $exportsLen }}{{ end }}
{{ "      ]," }}
      "export-memory-name": "memory"
    }
  }
}`

var mbtPluginPluginFunctionsTemplateStr = `{{range $index, $export := .Plugin.Exports }}{{ $name := .Name }}{{ if $index | lt 0 }}
{{ end }}/// Exported: {{ $name }}
pub fn exported_{{ $name | lowerSnakeCase }}() -> Int {
{{ if . | inputIsMbtVoidType }}  {{ $name | lowerSnakeCase }}(){{ end -}}
{{ if . | inputIsMbtPrimitiveType }}  let input = @host.input_string()
  let output = {{ $name | lowerSnakeCase }}(input) |> @jsonutil.to_json()
  @host.output_json_value(output){{ end -}}
{{ if . | inputIsMbtReferenceType }}  {{ inputMbtReferenceTypeName . }}::parse(@host.input_string())!!.unwrap()
  |> {{ $name | lowerSnakeCase }}()
  |> @jsonutil.to_json()
  |> @host.output_json_value(){{ end }}
  return 0 // success
{{ "}" }}
{{ end }}`
