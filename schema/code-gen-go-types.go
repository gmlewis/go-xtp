package schema

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

var (
	enumGoTemplate       = template.Must(template.New("enum.go").Funcs(funcMap).Parse(enumGoTemplateStr))
	structGoTemplate     = template.Must(template.New("struct.go").Funcs(funcMap).Parse(structGoTemplateStr))
	structTestGoTemplate = template.Must(template.New("struct_test.go").Funcs(funcMap).Parse(structTestGoTemplateStr))
)

// genGoCustomTypes generates custom types with tests for the plugin in Go.
func (p *Plugin) genGoCustomTypes() (srcOut, testSrcOut string, err error) {
	srcBlocks, testBlocks := make([]string, 0, len(p.CustomTypes)), make([]string, 0, len(p.CustomTypes))

	for _, ct := range p.CustomTypes {
		srcBlock, err := p.genGoCustomType(ct)
		if err != nil {
			return "", "", err
		}
		srcBlocks = append(srcBlocks, srcBlock)

		testBlock, err := p.genTestGoCustomType(ct)
		if err != nil {
			return "", "", err
		}
		testBlocks = append(testBlocks, testBlock)
	}

	srcToFmt := strings.Join(srcBlocks, "\n")
	src, err := format.Source([]byte(srcToFmt))
	if err != nil {
		return "", "", fmt.Errorf("gofmt error: %v\npre-formatted source:\n%v", err, srcToFmt)
	}

	testSrcToFmt := testPrelude + strings.Join(testBlocks, "\n")
	testSrc, err := format.Source([]byte(testSrcToFmt))
	if err != nil {
		return "", "", fmt.Errorf("gofmt error: %v\npre-formatted test source:\n%v", err, testSrcToFmt)
	}

	return string(src), string(testSrc), nil
}

// genGoCustomType generates Go source code for a single custom datatype.
func (p *Plugin) genGoCustomType(ct *CustomType) (string, error) {
	if ct == nil {
		return "", errors.New("unexpected nil CustomType")
	}

	switch {
	case len(ct.Enum) > 0:
		return p.genGoEnum(ct)
	case len(ct.Properties) > 0:
		return p.genGoStruct(ct)
	default:
		return "", fmt.Errorf("unhandled CustomType: %#v", *ct)
	}
}

// genTestGoCustomType generates Go source code for a single custom datatype.
func (p *Plugin) genTestGoCustomType(ct *CustomType) (string, error) {
	if ct == nil {
		return "", errors.New("unexpected nil CustomType")
	}

	switch {
	case len(ct.Enum) > 0:
		return "", nil // no enum tests written yet... possibly not necessary.
	case len(ct.Properties) > 0:
		return p.genTestGoStruct(ct)
	default:
		return "", fmt.Errorf("unhandled CustomType: %#v", *ct)
	}
}

// getGoEnum generates Go source code for a single enum custom datatype.
func (p *Plugin) genGoEnum(ct *CustomType) (string, error) {
	var buf bytes.Buffer
	if err := enumGoTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var enumGoTemplateStr = `{{ $name := .Name }}// {{ $name }} represents {{ .Description | downcaseFirst | multilineComment }}.
type {{ $name }} string

const (
  {{range .Enum}}{{ $name }}Enum{{ . | uppercaseFirst }} {{ $name }} = "{{ . }}"
  {{ end }}
)
`

// getGoStruct generates Go source code for a single struct custom datatype.
func (p *Plugin) genGoStruct(ct *CustomType) (string, error) {
	var buf bytes.Buffer
	if err := structGoTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getTestGoStruct generates Go source code for a single struct custom datatype.
func (p *Plugin) genTestGoStruct(ct *CustomType) (string, error) {
	var buf bytes.Buffer
	if err := structTestGoTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var structGoTemplateStr = `{{ $name := .Name }}{{ $top := . }}// {{ $name }} represents {{ .Description | downcaseFirst }}.
type {{ $name }} struct {
  {{range .Properties}}{{ .Description | optionalMultilineComment }}{{ .Name | uppercaseFirst }} {{ getGoType . }} ` + "`" + `json:"{{ .Name }}{{ addOmitIfNeeded . }}"` + "`" + `
  {{ end }}
}
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
{{range $index, $prop := .Properties}}{{if .IsRequired}}  {{ .Name | uppercaseFirst }}: {{ requiredValue . $top }},
{{ end }}{{ end }}
			},
			want: ` + "`" + `{
{{range $index, $prop := .Properties}}{{if .IsRequired}}  "{{ .Name }}": {{ requiredJSONValue . $top }}{{ showJSONCommaForRequired $index $top }}{{ end }}{{ end }}
}` + "`" + `,
		},
		{
			name: "optional fields",
			obj: &{{ .Name }}{
{{range $index, $prop := .Properties}}{{ if .IsRequired | not }}  {{ .Name | uppercaseFirst }}: {{ defaultValue . }},{{ end }}{{ end }}
			},
			want: ` + "`" + `{
{{range $index, $prop := .Properties}}  "{{ .Name }}": {{ defaultJSONValue . }}{{ showJSONCommaForOptional $index $top }}{{ end }}
}` + "`" + `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsoncomp.MarshalIndent(tt.obj, "", "  ")
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.want, string(got)); diff != "" {
				t.Log(string(got))
				t.Errorf("Marshal mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
`

var testPrelude = `import (
  "testing"

	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

func stringPtr(s string) *string { return &s }
`
