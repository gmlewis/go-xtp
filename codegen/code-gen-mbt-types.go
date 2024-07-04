package codegen

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/gmlewis/go-xtp/schema"
)

var (
	enumMbtTemplate       = template.Must(template.New("code-gen-mbt-types.go:enumMbtTemplateStr").Funcs(funcMap).Parse(enumMbtTemplateStr))
	structMbtTemplate     = template.Must(template.New("code-gen-mbt-types.go:structMbtTemplateStr").Funcs(funcMap).Parse(structMbtTemplateStr))
	structTestMbtTemplate = template.Must(template.New("code-gen-mbt-types.go:structTestMbtTemplateStr").Funcs(funcMap).Parse(structTestMbtTemplateStr))
)

// genMbtCustomTypes generates custom types with tests for the plugin in Go.
func (c *Client) genMbtCustomTypes() error {
	srcBlocks, testBlocks := make([]string, 0, len(c.Plugin.CustomTypes)+1), make([]string, 0, len(c.Plugin.CustomTypes))

	for _, ct := range c.Plugin.CustomTypes {
		srcBlock, err := c.genMbtCustomType(ct)
		if err != nil {
			return err
		}
		srcBlocks = append(srcBlocks, srcBlock)

		testBlock, err := c.genTestMbtCustomType(ct)
		if err != nil {
			return err
		}
		if testBlock != "" {
			testBlocks = append(testBlocks, testBlock)
		}
	}

	if c.numStructs > 0 {
		srcBlocks = append(srcBlocks, mbtXTPSchemaMap)
	}

	src := strings.Join(srcBlocks, "\n")
	c.CustTypesFilename = fmt.Sprintf("%v.%v", c.PkgName, c.Lang)
	c.CustTypes = src
	testSrc := strings.Join(testBlocks, "\n")
	c.CustTypesTestsFilename = fmt.Sprintf("%v_test.%v", c.PkgName, c.Lang)
	c.CustTypesTests = testSrc

	return nil
}

// genMbtCustomType generates MoonBit source code for a single custom datatype.
func (c *Client) genMbtCustomType(ct *schema.CustomType) (string, error) {
	if ct == nil {
		return "", errors.New("unexpected nil CustomType")
	}

	switch {
	case len(ct.Enum) > 0:
		return c.genMbtEnum(ct)
	case len(ct.Properties) > 0:
		c.numStructs++
		return c.genMbtStruct(ct)
	default:
		return "", fmt.Errorf("unhandled CustomType: %#v", *ct)
	}
}

// genTestMbtCustomType generates MoonBit source code for a single custom datatype.
func (c *Client) genTestMbtCustomType(ct *schema.CustomType) (string, error) {
	if ct == nil {
		return "", errors.New("unexpected nil CustomType")
	}

	switch {
	case len(ct.Enum) > 0:
		return "", nil // no enum tests written yet... possibly not necessary.
	case len(ct.Properties) > 0:
		return c.genTestMbtStruct(ct)
	default:
		return "", fmt.Errorf("unhandled CustomType: %#v", *ct)
	}
}

// getGoEnum generates MoonBit source code for a single enum custom datatype.
func (c *Client) genMbtEnum(ct *schema.CustomType) (string, error) {
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
} derive(Debug, Eq)

pub fn {{ $name }}::parse(s : String) -> {{ $name }}!String {
  match s {
{{range .Enum}}    "{{ . }}" => {{ . | uppercaseFirst }}
{{ end -}}
{{ "    _ => {" }}
      raise "not a {{ $name }}: \(s)"
    }
  }
}

pub impl @jsonutil.ToJson for {{ $name }} with to_json(self) {
  match self {
  {{range .Enum}}  {{ . | uppercaseFirst }} => @jsonutil.to_json("{{ . }}")
  {{ end -}}
  }
}
`

// getGoStruct generates MoonBit source code for a single struct custom datatype.
func (c *Client) genMbtStruct(ct *schema.CustomType) (string, error) {
	var buf bytes.Buffer
	if err := structMbtTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getTestGoStruct generates MoonBit source code for a single struct custom datatype.
func (c *Client) genTestMbtStruct(ct *schema.CustomType) (string, error) {
	var buf bytes.Buffer
	if err := structTestMbtTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var mbtXTPSchemaMap = `// XTPSchema describes the values and types of an XTP object
// in a language-agnostic format.
type XTPSchema Map[String, String]
`

var structMbtTemplateStr = `{{ $name := .Name }}{{ $top := . }}/// ` + "`" + `{{ $name }}` + "`" + ` represents {{ .Description | downcaseFirst }}.
pub struct {{ $name }} {
{{range .Properties}}  {{ .Description | optionalMbtMultilineComment }}{{ .Name | lowerSnakeCase }} : {{ getMbtType . }}
{{ end -}}
} derive(Debug, Eq)

/// ` + "`" + `{{ $name }}::new` + "`" + ` returns a new struct with default values.
pub fn {{ $name }}::new() -> {{ $name }} {
  {
{{range .Properties}}    {{ .Name | lowerSnakeCase }}: {{ defaultMbtValue . }},
{{ end -}}
{{ "  }" }}
}

pub impl @jsonutil.ToJson for {{ $name }} with to_json(self) {
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

/// ` + "`{{ $name }}::from_json` transforms a `@json.JsonValue`" + ` to a value.
pub fn {{ $name }}::from_json(value : @json.JsonValue) -> {{ $name }}? {
  match value {
    @json.JsonValue::Object({
{{range .Properties}}{{ if .IsRequired }}      "{{ .Name }}": Some(@json.JsonValue::String({{ .Name | lowerSnakeCase }})),
{{ end }}{{ if .IsRequired | not }}      "{{ .Name }}": {{ .Name | lowerSnakeCase }},
{{ end }}{{ end -}}
{{ "    }) => Some({" }}
{{range .Properties}}{{ if .IsRequired }}      {{ .Name | lowerSnakeCase }}: {{ .Name | lowerSnakeCase }},
{{ end }}{{ if .IsRequired | not }}      {{ .Name | lowerSnakeCase }}: match {{ .Name | lowerSnakeCase }} {
        Some({{ mbtFromJSONMatchKey . }}) => {{ mbtFromJSONMatchValue . }}
        _ => None
      },
{{ end }}{{ end -}}
{{ "    })" }}
    _ => None
  }
}

/// ` + "`{{ $name }}::parse` parses a JSON string and returns the value." + `
pub fn {{ $name }}::parse(s : String) -> {{ $name }}!String {
  match @json.parse(s) {
    Ok(jv) =>
      match {{ $name }}::from_json(jv) {
        Some(value) => value
        None => {
          raise "unable to parse {{ $name }} \(s)"
        }
      }
    Err(e) => {
      raise "unable to parse {{ $name }} \(s): \(e)"
    }
  }
}

/// ` + "`get_schema` returns an `XTPSchema` for the `{{ $name }}`" + `.
pub fn get_schema(self : {{ $name }}) -> XTPSchema {
  {
{{range .Properties}}    "{{ .Name }}": "{{ getExtismType . $top }}",
{{ end -}}
{{ "  }" }}
}
`

var structTestMbtTemplateStr = `{{ $name := .Name }}{{ $top := . }}test "{{ $name }}" {
  let default_object = {{ $name }}::new()
  let got = @jsonutil.to_json(default_object)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
{{ "    #|{" }}{{range $index, $prop := .Properties}}{{ if .IsRequired }}"{{ .Name }}":{{ defaultMbtJSONValue . $top }}{{ showJSONCommaForRequired $index $top }}{{ end }}{{ end -}}{{ "}" }}
  @test.eq(got, want)!
  //
  let got_parse = {{ $name }}::parse(want)!
  @test.eq(got_parse, default_object)!
  //
  let required_fields : {{ $name }} = {
{{range .Properties}}    {{ .Name | lowerSnakeCase }}: {{ requiredMbtValue . }},
{{ end -}}
{{ "  }" }}
  let got = @jsonutil.to_json(required_fields)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
{{ "    #|{" }}{{range $index, $prop := .Properties}}{{ if .IsRequired }}"{{ .Name }}":{{ requiredMbtJSONValue . $top }}{{ showJSONCommaForRequired $index $top }}{{ end }}{{ end -}}{{ "}" }}
  @test.eq(got, want)!
  //
  let got_parse = {{ $name }}::parse(want)!
  @test.eq(got_parse, required_fields)!
{{ if hasOptionalFields .}}  //
  let optional_fields : {{ $name }} = {
    ..required_fields,
{{range $index, $prop := .Properties}}{{ if .IsRequired | not}}    {{ .Name | lowerSnakeCase }}: {{ optionalMbtValue . }},
{{ end }}{{ end -}}
{{ "  }" }}
  let got = @jsonutil.to_json(optional_fields)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
{{ "    #|{" }}{{ $propLen := .Properties | len }}{{range $index, $prop := .Properties}}"{{ .Name }}":{{ requiredMbtJSONValue . $top }}{{ showJSONCommaForOptional $index $propLen }}{{ end -}}{{ "}" }}
  @test.eq(got, want)!
  //
  let got_parse = {{ $name }}::parse(want)!
  @test.eq(got_parse, optional_fields)!
{{ end -}}
}
`
