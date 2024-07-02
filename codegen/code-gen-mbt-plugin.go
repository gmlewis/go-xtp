package codegen

import (
	"bytes"
	"text/template"
)

var (
	mbtPluginXtpTOMLTemplate       = template.Must(template.New("cost-gen-mbt-plugin.go:mbtPluginXtpTOMLTemplateStr").Parse(mbtPluginXtpTOMLTemplateStr))
	mbtPluginHostFunctionsTemplate = template.Must(template.New("cost-gen-mbt-plugin.go:mbtPluginHostFunctionsTemplateStr").Funcs(funcMap).Parse(mbtPluginHostFunctionsTemplateStr))
	mbtPluginMainTemplate          = template.Must(template.New("cost-gen-mbt-plugin.go:mbtPluginMainTemplateStr").Funcs(funcMap).Parse(mbtPluginMainTemplateStr))
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

	m := GeneratedFiles{
		"build.sh":               buildShScript,
		c.CustTypesFilename:      c.CustTypes,
		c.CustTypesTestsFilename: c.CustTypesTests,
		"host-functions.mbt":     hostFunctionsStr.String(),
		"main.mbt":               mainStr.String(),
		"moon.pkg.json":          "", // TODO
		"plugin-functions.mbt":   "", // TODO
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

var mbtPluginHostFunctionsTemplateStr = `{{ $top := . }}{{range .Plugin.Imports }}{{ $name := .Name }}pub fn host_{{ $name | lowerSnakeCase }}(offset : Int64) -> Int64 = "extism:host/user" "{{ $name }}"

/// ` + "`{{ $name | lowerSnakeCase }}`" + ` - {{ .Description | mbtMultilineComment | stripLeadingSlashes | leftJustify }}
pub fn {{ $name | lowerSnakeCase }}(input : {{ .Input | inputToMbtType }}) -> {{ .Output | outputToMbtType }}!String {
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

var mbtPluginMainTemplateStr = `{{ $top := . }}{{range .Plugin.Exports }}{{ $name := .Name }}/// ` + "`{{ $name | lowerSnakeCase }}`" + ` - {{ .Description | mbtMultilineComment | stripLeadingSlashes | leftJustify }}
pub fn {{ $name | lowerSnakeCase }}({{ .Input | inputToMbtType }}) -> {{ .Output | outputToMbtType }} {
  // TODO: fill out your implementation here
  return-type-here
}

{{ end }}fn main {

}
`
