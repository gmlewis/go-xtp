package codegen

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gmlewis/go-xtp/schema"
)

var (
	unreachable = errors.New("unreachable")
)

func (c *Client) GenCustomTypes() (GeneratedFiles, error) {
	m := GeneratedFiles{
		c.CustTypesFilename:      c.CustTypes,
		c.CustTypesTestsFilename: c.CustTypesTests,
	}

	switch c.Lang {
	case "go":
		m[c.CustTypesFilename] = fmt.Sprintf("// Package %v represents the custom datatypes for an XTP Extension Plugin.\npackage %[1]v\n\n%v", c.PkgName, c.CustTypes)
		m[c.CustTypesTestsFilename] = fmt.Sprintf("package %v\n\n%v", c.PkgName, c.CustTypesTests)
	case "mbt":
		m["moon.pkg.json"] = defaultMoonPkgJSONFile
	}

	return m, nil
}

// genCustomTypes generates custom types with tests for the plugin.
func (c *Client) genCustomTypes() error {
	switch c.Lang {
	case "go":
		return c.genGoCustomTypes()
	case "mbt":
		return c.genMbtCustomTypes()
	}
	return unreachable
}

// GenHostSDK generates Host SDK code to call the extension plugin.
func (c *Client) GenHostSDK() (GeneratedFiles, error) {
	switch c.Lang {
	case "go":
		return c.genGoHostSDK()
	case "mbt":
		return c.genMbtHostSDK()
	}
	return GeneratedFiles{}, unreachable
}

// GenPluginPDK generates Plugin PDK code to process plugin calls.
func (c *Client) GenPluginPDK() (GeneratedFiles, error) {
	switch c.Lang {
	case "go":
		return c.genGoPluginPDK()
	case "mbt":
		return c.genMbtPluginPDK()
	}
	return GeneratedFiles{}, unreachable
}

var funcMap = map[string]any{
	"addOmitIfNeeded":             addOmitIfNeeded,
	"defaultGoJSONValue":          defaultGoJSONValue,
	"defaultGoValue":              defaultGoValue,
	"defaultMbtJSONValue":         defaultMbtJSONValue,
	"defaultMbtValue":             defaultMbtValue,
	"downcaseFirst":               downcaseFirst,
	"getExtismType":               getExtismType,
	"getGoType":                   getGoType,
	"getMbtType":                  getMbtType,
	"hasOptionalFields":           hasOptionalFields,
	"inputToMbtType":              inputToMbtType,
	"jsonOutputAsMbtType":         jsonOutputAsMbtType,
	"leftJustify":                 leftJustify,
	"lowerSnakeCase":              lowerSnakeCase,
	"mbtMultilineComment":         mbtMultilineComment,
	"multilineComment":            multilineComment,
	"optionalGoMultilineComment":  optionalGoMultilineComment,
	"optionalMbtMultilineComment": optionalMbtMultilineComment,
	"optionalMbtValue":            optionalMbtValue,
	"outputToMbtType":             outputToMbtType,
	"requiredGoJSONValue":         requiredGoJSONValue,
	"requiredGoValue":             requiredGoValue,
	"requiredMbtJSONValue":        requiredMbtJSONValue,
	"requiredMbtValue":            requiredMbtValue,
	"showJSONCommaForOptional":    showJSONCommaForOptional,
	"showJSONCommaForRequired":    showJSONCommaForRequired,
	"stripLeadingSlashes":         stripLeadingSlashes,
	"uppercaseFirst":              uppercaseFirst,
}

func addOmitIfNeeded(prop *schema.Property) string {
	if prop.IsRequired {
		return ""
	}
	return ",omitempty"
}

func downcaseFirst(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	return strings.ToLower(s[0:1]) + s[1:]
}

func getExtismType(prop *schema.Property, ct *schema.CustomType) string {
	var optionalMark string
	if !prop.IsRequired {
		optionalMark = "?"
	}

	var extismType string
	switch prop.Type {
	case "integer", "number", "boolean":
		extismType = prop.Type
	case "string":
		extismType = prop.Type
		if prop.Format == "date-time" {
			extismType = "Date"
		}
	default:
		if prop.Ref != "" {
			parts := strings.Split(prop.Ref, "/")
			extismType = parts[len(parts)-1]
		} else {
			log.Printf("WARNING: unknown property type %q: %#v", prop.Type, prop)
		}
	}

	return optionalMark + extismType
}

func hasOptionalFields(ct *schema.CustomType) bool {
	return len(ct.Required) != len(ct.Properties)
}

func leftJustify(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}

func lowerSnakeCase(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	result := strings.ToLower(s[0:1])
	for _, r := range s[1:] {
		rs := string(r)
		if rs == strings.ToUpper(rs) {
			result += "_" + strings.ToLower(rs)
			continue
		}
		result += rs
	}
	return result
}

func multilineComment(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n// ")
}

func showJSONCommaForOptional(index int, ct *schema.CustomType) string {
	if index < len(ct.Properties)-1 {
		return ","
	}
	return ""
}

func showJSONCommaForRequired(index int, ct *schema.CustomType) string {
	for index++; index < len(ct.Properties); index++ {
		prop := ct.Properties[index]
		if prop.IsRequired {
			return ","
		}
	}
	return ""
}

func stripLeadingSlashes(s string) string {
	return strings.TrimLeft(s, "/ ")
}

func uppercaseFirst(s string) string {
	if len(s) < 2 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

const defaultMoonPkgJSONFile = `{
  "import": [
    {
      "path": "gmlewis/json",
      "alias": "jsonutil"
    }
  ]
}`
