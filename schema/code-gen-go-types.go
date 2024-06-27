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

// getGoEnum generates Go source code for a single enum custom datatype.
func (p *Plugin) genGoEnum(ct *CustomType) (string, error) {
	t := template.Must(template.New("enum-go").Funcs(funcMap).Parse(enumGoTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, ct); err != nil {
		return "", err
	}
	out, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(out), nil
}

var enumGoTemplate = `{{ $name := .Name }}// {{ $name }} represents {{ .Description | downcaseFirst | multilineComment }}.
type {{ $name }} string

const (
  {{range .Enum}}{{ $name }}Enum{{ . | uppercaseFirst }} {{ $name }} = "{{ . }}"
  {{ end }}
)
`

// getGoStruct generates Go source code for a single struct custom datatype.
func (p *Plugin) genGoStruct(ct *CustomType) (string, error) {
	t := template.Must(template.New("struct-go").Funcs(funcMap).Parse(structGoTemplate))
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
