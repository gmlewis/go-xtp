package codegen

import (
	"errors"
	"strings"

	"github.com/gmlewis/go-xtp/schema"
)

var (
	unreachable = errors.New("unreachable")
)

// genCustomTypes generates custom types with tests for the plugin.
func (c *Client) genCustomTypes() (srcFile, testFile string, err error) {
	switch c.Lang {
	case "go":
		return c.genGoCustomTypes()
	case "mbt":
		return c.genMbtCustomTypes()
	}
	return "", "", unreachable
}

// GenHostSDK generates Host SDK code to call the extension plugin.
func (c *Client) GenHostSDK() (string, error) {
	switch c.Lang {
	case "go":
		return c.genGoHostSDK()
	case "mbt":
		return c.genMbtHostSDK()
	}
	return "", unreachable
}

// GenPluginPDK generates Plugin PDK code to process plugin calls.
func (c *Client) GenPluginPDK() (string, error) {
	switch c.Lang {
	case "go":
		return c.genGoPluginPDK()
	case "mbt":
		return c.genMbtPluginPDK()
	}
	return "", unreachable
}

var funcMap = map[string]any{
	"addOmitIfNeeded":             addOmitIfNeeded,
	"defaultGoJSONValue":          defaultGoJSONValue,
	"defaultGoValue":              defaultGoValue,
	"defaultMbtJSONValue":         defaultMbtJSONValue,
	"defaultMbtValue":             defaultMbtValue,
	"downcaseFirst":               downcaseFirst,
	"getGoType":                   getGoType,
	"getMbtType":                  getMbtType,
	"hasOptionalFields":           hasOptionalFields,
	"lowerSnakeCase":              lowerSnakeCase,
	"multilineComment":            multilineComment,
	"optionalGoMultilineComment":  optionalGoMultilineComment,
	"optionalMbtMultilineComment": optionalMbtMultilineComment,
	"optionalMbtValue":            optionalMbtValue,
	"requiredGoJSONValue":         requiredGoJSONValue,
	"requiredGoValue":             requiredGoValue,
	"requiredMbtJSONValue":        requiredMbtJSONValue,
	"requiredMbtValue":            requiredMbtValue,
	"showJSONCommaForOptional":    showJSONCommaForOptional,
	"showJSONCommaForRequired":    showJSONCommaForRequired,
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

func hasOptionalFields(ct *schema.CustomType) bool {
	return len(ct.Required) != len(ct.Properties)
}

func lowerSnakeCase(s string) string {
	if s == "" {
		return s
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

func uppercaseFirst(s string) string {
	if len(s) < 2 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}
