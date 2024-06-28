package schema

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoCodeGeneration = errors.New("code generation not supported for version v0")
)

// GenCustomTypes generates custom types with tests for the plugin.
func (p *Plugin) GenCustomTypes() (srcFile, testFile string, err error) {
	if p.Version == "v0" {
		return "", "", ErrNoCodeGeneration
	}

	switch p.Lang {
	case "go":
		return p.genGoCustomTypes()
	case "mbt":
		return p.genMbtCustomTypes()
	default:
		return "", "", fmt.Errorf("unsupported programming language: %q", p.Lang)
	}
}

// GenHostSDK generates Host SDK code to call the extension plugin.
func (p *Plugin) GenHostSDK(customTypes, pkgName string) (string, error) {
	if p.Version == "v0" {
		return "", ErrNoCodeGeneration
	}

	switch p.Lang {
	case "go":
		return p.genGoHostSDK(customTypes, pkgName)
	case "mbt":
		return p.genMbtHostSDK(customTypes, pkgName)
	default:
		return "", fmt.Errorf("unsupported programming language: %q", p.Lang)
	}
}

// GenPluginPDK generates Plugin PDK code to process plugin calls.
func (p *Plugin) GenPluginPDK(customTypes, pkgName string) (string, error) {
	if p.Version == "v0" {
		return "", ErrNoCodeGeneration
	}

	switch p.Lang {
	case "go":
		return p.genGoPluginPDK(customTypes, pkgName)
	case "mbt":
		return p.genMbtPluginPDK(customTypes, pkgName)
	default:
		return "", fmt.Errorf("unsupported programming language: %q", p.Lang)
	}
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

func addOmitIfNeeded(prop *Property) string {
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

func hasOptionalFields(ct *CustomType) bool {
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

func showJSONCommaForOptional(index int, ct *CustomType) string {
	if index < len(ct.Properties)-1 {
		return ","
	}
	return ""
}

func showJSONCommaForRequired(index int, ct *CustomType) string {
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
