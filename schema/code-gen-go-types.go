package schema

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

// genGoCustomTypes generates custom types with tests for the plugin in Go.
func (p *Plugin) genGoCustomTypes() (srcFile, testFile string, err error) {
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

	return strings.Join(srcBlocks, "\n"), strings.Join(testBlocks, "\n"), nil
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
		return p.genTestGoEnum(ct)
	case len(ct.Properties) > 0:
		return p.genTestGoStruct(ct)
	default:
		return "", fmt.Errorf("unhandled CustomType: %#v", *ct)
	}
}

// getGoEnum generates Go source code for a single enum custom datatype.
func (p *Plugin) genGoEnum(ct *CustomType) (string, error) {
	t := template.Must(template.New("enum.go").Funcs(funcMap).Parse(enumGoTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, ct); err != nil {
		return "", err
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(src), nil
}

// getTestGoEnum generates Go source code for a single enum custom datatype.
func (p *Plugin) genTestGoEnum(ct *CustomType) (string, error) {
	t := template.Must(template.New("enum_test.go").Funcs(funcMap).Parse(enumTestGoTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, ct); err != nil {
		return "", err
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(src), nil
}

var enumGoTemplate = `{{ $name := .Name }}// {{ $name }} represents {{ .Description | downcaseFirst | multilineComment }}.
type {{ $name }} string

const (
  {{range .Enum}}{{ $name }}Enum{{ . | uppercaseFirst }} {{ $name }} = "{{ . }}"
  {{ end }}
)
`

var enumTestGoTemplate = `{{ $name := .Name }}// {{ $name }} represents {{ .Description | downcaseFirst | multilineComment }}.
type {{ $name }} string

const (
  {{range .Enum}}{{ $name }}Enum{{ . | uppercaseFirst }} {{ $name }} = "{{ . }}"
  {{ end }}
)
`

// getGoStruct generates Go source code for a single struct custom datatype.
func (p *Plugin) genGoStruct(ct *CustomType) (string, error) {
	t := template.Must(template.New("struct.go").Funcs(funcMap).Parse(structGoTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, ct); err != nil {
		return "", err
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("gofmt error: %v\npre-formatted source:\n%v", err, buf.String())
	}

	return string(out), nil
}

// getTestGoStruct generates Go source code for a single struct custom datatype.
func (p *Plugin) genTestGoStruct(ct *CustomType) (string, error) {
	t := template.Must(template.New("struct_test.go").Funcs(funcMap).Parse(structTestGoTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, ct); err != nil {
		return "", err
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("gofmt error: %v\npre-formatted source:\n%v", err, buf.String())
	}

	return string(out), nil
}

var structGoTemplate = `{{ $name := .Name }}{{ $top := . }}// {{ $name }} represents {{ .Description | downcaseFirst }}.
type {{ $name }} struct {
  {{range .Properties}}// {{ .Description | multilineComment }}
  {{ .Name | uppercaseFirst }} {{ getGoType . }} ` + "`" + `json:"{{ .Name }}{{ addOmitIfNeeded . }}"` + "`" + `
  {{ end }}
}
`

var structTestGoTemplate = `{{ $name := .Name }}{{ $top := . }}import (
  "testing"

	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

func stringPtr(s string) *string { return &s }

func Test{{ $name }}Marshal(t *testing.T) {
  t.Parallel()
	tests := []struct {
		name string
		obj  *{{ .Name }}
		want string
	}{
		{
			name: "required fields",
			obj: &{{ .Name }}{
				Ghost:    GhostGangEnumBlinky,
				ABoolean: true,
				AString:  "aString",
				AnInt:    0,
			},
			want: ` + "`" + `{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0}` + "`" + `,
		},
		{
			name: "optional fields",
			obj: &{{ .Name }}{
				AnOptionalDate: stringPtr("anOptionalDate"),
			},
			want: ` + "`" + `{"ghost":"","aBoolean":false,"aString":"","anInt":0,"anOptionalDate":"anOptionalDate"}` + "`" + `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsoncomp.Marshal(tt.obj)
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
