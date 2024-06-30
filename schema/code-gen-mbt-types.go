package schema

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

var (
	enumMbtTemplate       = template.Must(template.New("code-gen-mbt-types.go:enumMbtTemplateStr").Funcs(funcMap).Parse(enumMbtTemplateStr))
	structMbtTemplate     = template.Must(template.New("code-gen-mbt-types.go:structMbtTemplateStr").Funcs(funcMap).Parse(structMbtTemplateStr))
	structTestMbtTemplate = template.Must(template.New("code-gen-mbt-types.go:structTestMbtTemplateStr").Funcs(funcMap).Parse(structTestMbtTemplateStr))
)

// genMbtCustomTypes generates custom types with tests for the plugin in Go.
func (p *Plugin) genMbtCustomTypes() (srcOut, testSrcOut string, err error) {
	srcBlocks, testBlocks := make([]string, 0, len(p.CustomTypes)), make([]string, 0, len(p.CustomTypes))

	for _, ct := range p.CustomTypes {
		srcBlock, err := p.genMbtCustomType(ct)
		if err != nil {
			return "", "", err
		}
		srcBlocks = append(srcBlocks, srcBlock)

		testBlock, err := p.genTestMbtCustomType(ct)
		if err != nil {
			return "", "", err
		}
		if testBlock != "" {
			testBlocks = append(testBlocks, testBlock)
		}
	}

	src := strings.Join(srcBlocks, "\n")
	testSrc := strings.Join(testBlocks, "\n")

	return string(src), string(testSrc), nil
}

// genMbtCustomType generates MoonBit source code for a single custom datatype.
func (p *Plugin) genMbtCustomType(ct *CustomType) (string, error) {
	if ct == nil {
		return "", errors.New("unexpected nil CustomType")
	}

	switch {
	case len(ct.Enum) > 0:
		return p.genMbtEnum(ct)
	case len(ct.Properties) > 0:
		return p.genMbtStruct(ct)
	default:
		return "", fmt.Errorf("unhandled CustomType: %#v", *ct)
	}
}

// genTestMbtCustomType generates MoonBit source code for a single custom datatype.
func (p *Plugin) genTestMbtCustomType(ct *CustomType) (string, error) {
	if ct == nil {
		return "", errors.New("unexpected nil CustomType")
	}

	switch {
	case len(ct.Enum) > 0:
		return "", nil // no enum tests written yet... possibly not necessary.
	case len(ct.Properties) > 0:
		return p.genTestMbtStruct(ct)
	default:
		return "", fmt.Errorf("unhandled CustomType: %#v", *ct)
	}
}

// getGoEnum generates MoonBit source code for a single enum custom datatype.
func (p *Plugin) genMbtEnum(ct *CustomType) (string, error) {
	var buf bytes.Buffer
	if err := enumMbtTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var enumMbtTemplateStr = `{{ $name := .Name }}/// ` + "`" + `{{ $name }}` + "`" + ` represents {{ .Description | downcaseFirst | multilineComment }}.
pub enum {{ $name }} {
{{range .Enum}}  {{ . | uppercaseFirst }}
{{ end -}}
}

impl @jsonutil.ToJson for {{ $name }} with to_json(self) {
  match self {
  {{range .Enum}}  {{ . | uppercaseFirst }} => @jsonutil.to_json("{{ . }}")
  {{ end -}}
  }
}
`

// getGoStruct generates MoonBit source code for a single struct custom datatype.
func (p *Plugin) genMbtStruct(ct *CustomType) (string, error) {
	var buf bytes.Buffer
	if err := structMbtTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getTestGoStruct generates MoonBit source code for a single struct custom datatype.
func (p *Plugin) genTestMbtStruct(ct *CustomType) (string, error) {
	var buf bytes.Buffer
	if err := structTestMbtTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var structMbtTemplateStr = `{{ $name := .Name }}{{ $top := . }}/// ` + "`" + `{{ $name }}` + "`" + ` represents {{ .Description | downcaseFirst }}.
pub struct {{ $name }} {
{{range .Properties}}  {{ .Description | optionalMbtMultilineComment }}{{ .Name | lowerSnakeCase }} : {{ getMbtType . }}
{{ end -}}
}

/// ` + "`" + `{{ $name }}::new` + "`" + ` returns a new struct with default values.
pub fn {{ $name }}::new() -> {{ $name }} {
  {
{{range .Properties}}    {{ .Name | lowerSnakeCase }}: {{ defaultMbtValue . }},
{{ end -}}
{{ "  }" }}
}

impl @jsonutil.ToJson for {{ $name }} with to_json(self) {
  let fields : Array[(String, @jsonutil.ToJson)] = [
{{range .Properties}}{{ if .IsRequired }}    ("{{ .Name }}", self.{{ .Name | lowerSnakeCase }}),
{{ end }}{{ end -}}
{{ "  ]" }}
{{range .Properties}}{{ if .IsRequired | not }}  match self.{{ .Name | lowerSnakeCase }} {
    Some(value) => fields.append([("{{ .Name }}", value)])
    None => ()
  }
{{ end }}{{ end -}}
{{ "  @jsonutil.from_entries(fields)" }}
}
`

var structTestMbtTemplateStr = `{{ $name := .Name }}{{ $top := . }}test "{{ $name }}" {
  let default_object = {{ $name }}::new()
  let got = default_object.to_json()
  let want =
{{ "    #|{" }}{{range $index, $prop := .Properties}}{{ if .IsRequired }}"{{ .Name }}":{{ defaultMbtJSONValue . $top }}{{ showJSONCommaForRequired $index $top }}{{ end }}{{ end -}}{{ "}" }}
  @assertion.assert_eq(got, want)?
  //
  let required_fields : {{ $name }} = {
{{range .Properties}}    {{ .Name | lowerSnakeCase }}: {{ requiredMbtValue . }},
{{ end -}}
{{ "  }" }}
  let got = required_fields.to_json()
  let want =
{{ "    #|{" }}{{range $index, $prop := .Properties}}{{ if .IsRequired }}"{{ .Name }}":{{ requiredMbtJSONValue . $top }}{{ showJSONCommaForRequired $index $top }}{{ end }}{{ end -}}{{ "}" }}
  @assertion.assert_eq(got, want)?
{{ if hasOptionalFields .}}  //
  let optional_fields : {{ $name }} = {
    ..required_fields,
{{range $index, $prop := .Properties}}{{ if .IsRequired | not}}    {{ .Name | lowerSnakeCase }}: {{ optionalMbtValue . }},
{{ end }}{{ end -}}
{{ "  }" }}
  let got = optional_fields.to_json()
  let want =
{{ "    #|{" }}{{range $index, $prop := .Properties}}"{{ .Name }}":{{ requiredMbtJSONValue . $top }}{{ showJSONCommaForOptional $index $top }}{{ end -}}{{ "}" }}
  @assertion.assert_eq(got, want)?
{{ end -}}
}
`
