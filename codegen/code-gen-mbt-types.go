package codegen

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/gmlewis/go-xtp/schema"
)

var (
	enumMbtTemplate       = template.Must(template.New("code-gen-mbt-types.go:enumMbtTemplateStr").Funcs(funcMap).Parse(enumMbtTemplateStr))
	enumTestMbtTemplate   = template.Must(template.New("code-gen-mbt-types.go:enumTestMbtTemplateStr").Funcs(funcMap).Parse(enumTestMbtTemplateStr))
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
	c.CustTypesTestsFilename = fmt.Sprintf("%v_bbtest.%v", c.PkgName, c.Lang)
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
		return c.getTestMbtEnum(ct)
	case len(ct.Properties) > 0:
		return c.genTestMbtStruct(ct)
	default:
		return "", fmt.Errorf("unhandled CustomType: %#v", *ct)
	}
}

// getMbtEnum generates MoonBit source code for a single enum custom datatype.
func (c *Client) genMbtEnum(ct *schema.CustomType) (string, error) {
	var buf bytes.Buffer
	if err := enumMbtTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getTestMbtEnum generates MoonBit test source code for a single enum custom datatype.
func (c *Client) getTestMbtEnum(ct *schema.CustomType) (string, error) {
	var buf bytes.Buffer
	if err := enumTestMbtTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

//go:embed enum-mbt-template.txt
var enumMbtTemplateStr string

//go:embed enum-test-mbt-template.txt
var enumTestMbtTemplateStr string

// getMbtStruct generates MoonBit source code for a single struct custom datatype.
func (c *Client) genMbtStruct(ct *schema.CustomType) (string, error) {
	var buf bytes.Buffer
	if err := structMbtTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getTestMbtStruct generates MoonBit test source code for a single struct custom datatype.
func (c *Client) genTestMbtStruct(ct *schema.CustomType) (string, error) {
	var buf bytes.Buffer
	if err := structTestMbtTemplate.Execute(&buf, ct); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var mbtXTPSchemaMap = "/// `XTPSchema` describes the values and types of an XTP object" + `
/// in a language-agnostic format.
type XTPSchema Map[String, String]
`

//go:embed struct-mbt-template.txt
var structMbtTemplateStr string

//go:embed struct-test-mbt-template.txt
var structTestMbtTemplateStr string
