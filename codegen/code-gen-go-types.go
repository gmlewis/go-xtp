package codegen

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/gmlewis/go-xtp/schema"
)

var (
	enumGoTemplate       = template.Must(template.New("code-gen-go-types.go:enumGoTemplateStr").Funcs(funcMap).Parse(enumGoTemplateStr))
	enumTestGoTemplate   = template.Must(template.New("code-gen-go-types.go:enumTestGoTemplateStr").Funcs(funcMap).Parse(enumTestGoTemplateStr))
	structGoTemplate     = template.Must(template.New("code-gen-go-types.go:structGoTemplateStr").Funcs(funcMap).Parse(structGoTemplateStr))
	structTestGoTemplate = template.Must(template.New("code-gen-go-types.go:structTestGoTemplateStr").Funcs(funcMap).Parse(structTestGoTemplateStr))
)

// genGoCustomTypes generates custom types with tests for the plugin in Go.
func (c *Client) genGoCustomTypes() error {
	srcBlocks, testBlocks := make([]string, 0, len(c.Plugin.CustomTypes)+1), make([]string, 0, len(c.Plugin.CustomTypes))

	for _, ct := range c.Plugin.CustomTypes {
		srcBlock, err := c.genGoCustomType(ct)
		if err != nil {
			return err
		}
		srcBlocks = append(srcBlocks, srcBlock)

		testBlock, err := c.genTestGoCustomType(ct)
		if err != nil {
			return err
		}
		testBlocks = append(testBlocks, testBlock)
	}

	if c.numStructs > 0 {
		srcBlocks = append(srcBlocks, goXTPSchemaMap)
	}

	srcToFmt := strings.Join(srcBlocks, "\n")
	if strings.Contains(srcToFmt, "fmt.Errorf") {
		srcToFmt = goPreludeWithFmt + srcToFmt
	} else {
		srcToFmt = goPrelude + srcToFmt
	}
	src, err := format.Source([]byte(srcToFmt))
	if err != nil {
		return fmt.Errorf("gofmt error: %v\npre-formatted source:\n%v", err, srcToFmt)
	}
	c.CustTypesFilename = fmt.Sprintf("%v.%v", c.PkgName, c.Lang)
	c.CustTypes = string(src)

	testSrcToFmt := testGoPrelude + strings.Join(testBlocks, "\n")
	testSrc, err := format.Source([]byte(testSrcToFmt))
	if err != nil {
		return fmt.Errorf("gofmt error: %v\npre-formatted test source:\n%v", err, testSrcToFmt)
	}
	c.CustTypesTestsFilename = fmt.Sprintf("%v_test.%v", c.PkgName, c.Lang)
	c.CustTypesTests = string(testSrc)

	return nil
}

// genGoCustomType generates Go source code for a single custom datatype.
func (c *Client) genGoCustomType(ct *schema.CustomType) (string, error) {
	if ct == nil {
		return "", errors.New("unexpected nil CustomType")
	}

	switch {
	case len(ct.Enum) > 0:
		return c.genGoEnum(ct)
	case len(ct.Properties) > 0:
		c.numStructs++
		return c.genGoStruct(ct)
	default:
		return "", fmt.Errorf("unhandled CustomType: %#v", *ct)
	}
}

// genTestGoCustomType generates Go test source code for a single custom datatype.
func (c *Client) genTestGoCustomType(ct *schema.CustomType) (string, error) {
	if ct == nil {
		return "", errors.New("unexpected nil CustomType")
	}

	switch {
	case len(ct.Enum) > 0:
		return c.genTestGoEnum(ct)
	case len(ct.Properties) > 0:
		return c.genTestGoStruct(ct)
	default:
		return "", fmt.Errorf("unhandled CustomType: %#v", *ct)
	}
}

// getGoEnum generates Go source code for a single enum custom datatype.
func (c *Client) genGoEnum(ct *schema.CustomType) (string, error) {
	var buf bytes.Buffer
	if err := enumGoTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getTestGoEnum generates Go test source code for a single enum custom datatype.
func (c *Client) genTestGoEnum(ct *schema.CustomType) (string, error) {
	var buf bytes.Buffer
	if err := enumTestGoTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var enumGoTemplateStr = `{{ $name := .Name }}// {{ $name }} represents {{ .Description | downcaseFirst | multilineComment }}.
type {{ $name }} string

const (
{{range .Enum}}  {{ $name }}Enum{{ . | uppercaseFirst }} {{ $name }} = "{{ . }}"
{{ end -}}
)

// Parse{{ $name }} parses a JSON string and returns the value.
func Parse{{ $name }}(s string) (value {{ $name }}, err error) {
	switch s {
` + "{{range .Enum}}	case `\"{{ . }}\"`:" + `
		return {{ $name }}Enum{{ . | uppercaseFirst }}, nil
{{ end -}}
	default:
		return value, fmt.Errorf("not a {{ $name }}: %v", s)
	}
}
`

var enumTestGoTemplateStr = `{{ $name := .Name }}{{ $top := . }}func TestParse{{ $name }}(t *testing.T) {
	t.Parallel()

	{{ $name | downcaseFirst }} := {{ $name }}Enum{{ index .Enum 0 | uppercaseFirst }}
	buf, err := jsoncomp.Marshal({{ $name | downcaseFirst }})
	if err != nil {
		t.Fatal(err)
	}

` + "	want := `\"{{ index .Enum 0 }}\"`" + `
	if got := string(buf); got != want {
		t.Errorf("Marshal = '%v', want '%v'", got, want)
	}

	got, err := Parse{{ $name }}(want)
	if err != nil {
		t.Fatal(err)
	}
	if got != {{ $name | downcaseFirst }} {
		t.Errorf("Parse{{ $name }} = '%v', want '%v'", got, {{ $name | downcaseFirst }})
	}
}
`

// getGoStruct generates Go source code for a single struct custom datatype.
func (c *Client) genGoStruct(ct *schema.CustomType) (string, error) {
	var buf bytes.Buffer
	if err := structGoTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getTestGoStruct generates Go test source code for a single struct custom datatype.
func (c *Client) genTestGoStruct(ct *schema.CustomType) (string, error) {
	var buf bytes.Buffer
	if err := structTestGoTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var goXTPSchemaMap = `// XTPSchema describes the values and types of an XTP object
// in a language-agnostic format.
type XTPSchema map[string]string
`

var structGoTemplateStr = `{{ $name := .Name }}{{ $top := . }}// {{ $name }} represents {{ .Description | downcaseFirst }}.
type {{ $name }} struct {
{{range .Properties}}  {{ .Description | optionalGoMultilineComment }}{{ .Name | uppercaseFirst }} {{ getGoType . }} ` + "`" + `json:"{{ .Name }}{{ addOmitIfNeeded . }}"` + "`" + `
{{ end -}}
}

// Parse{{ $name }} parses a JSON string and returns the value.
func Parse{{ $name }}(s string) (value {{ $name }}, err error) {
	if err := json.Unmarshal([]byte(s), &value); err != nil {
		return value, err
	}

	return value, nil
}

// GetSchema returns an ` + "`" + `XTPSchema` + "`" + ` for the ` + "`" + `{{ $name }}` + "`" + `.
func (c *{{ $name }}) GetSchema() XTPSchema {
	return XTPSchema{
{{range .Properties}}    "{{ .Name }}": "{{ getExtismType . $top }}",
{{ end -}}
	}
}
`

var goPrelude = `import (
	"encoding/json" // jsoniter/jsoncomp are not compatible with tinygo.
)

`

var goPreludeWithFmt = `import (
	"encoding/json" // jsoniter/jsoncomp are not compatible with tinygo.
	"fmt"
)

`

var testGoPrelude = `import (
	"testing"

	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

func boolPtr(b bool) *bool { return &b }
func intPtr(i int) *int { return &i }
func stringPtr(s string) *string { return &s }

`

var structTestGoTemplateStr = `{{ $name := .Name }}{{ $top := . }}func Test{{ $name }}Marshal(t *testing.T) {
  t.Parallel()
	tests := []struct {
		name string
		obj  *{{ .Name }}
		want string
	}{
		{
			name: "required fields",
			obj: &{{ .Name }}{
{{range $index, $prop := .Properties}}{{if .IsRequired}}  {{ .Name | uppercaseFirst }}: {{ requiredGoValue . }},
{{ end }}{{ end }}
			},
			want: ` + "`" + `{{"{"}}{{range $index, $prop := .Properties}}{{if .IsRequired}}"{{ .Name }}":{{ requiredGoJSONValue . }}{{ showJSONCommaForRequired $index $top }}{{ end }}{{ end }}{{"}"}}` + "`" + `,
		},
		{
			name: "optional fields",
			obj: &{{ .Name }}{
{{range $index, $prop := .Properties}}{{ if .IsRequired | not }}  {{ .Name | uppercaseFirst }}: {{ defaultGoValue . }},
{{ end }}{{ end }}
			},
			want: ` + "`" + `{{"{"}}{{ $propLen := .Properties | len }}{{range $index, $prop := .Properties}}"{{ .Name }}":{{ defaultGoJSONValue . $top }}{{ showJSONCommaForOptional $index $propLen }}{{ end }}{{"}"}}` + "`" + `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsoncomp.Marshal(tt.obj)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.want, string(got)); diff != "" {
				t.Logf("got:\n%v", string(got))
				t.Errorf("Marshal mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
`
