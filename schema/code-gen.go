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
	"addOmitIfNeeded":          addOmitIfNeeded,
	"defaultJSONValue":         defaultJSONValue,
	"defaultValue":             defaultValue,
	"downcaseFirst":            downcaseFirst,
	"getGoType":                getGoType,
	"multilineComment":         multilineComment,
	"optionalMultilineComment": optionalMultilineComment,
	"requiredJSONValue":        requiredJSONValue,
	"requiredValue":            requiredValue,
	"showJSONCommaForOptional": showJSONCommaForOptional,
	"showJSONCommaForRequired": showJSONCommaForRequired,
	"uppercaseFirst":           uppercaseFirst,
}

func addOmitIfNeeded(prop *Property) string {
	if prop.IsRequired {
		return ""
	}
	return ",omitempty"
}

func defaultJSONValue(prop *Property) string {
	switch prop.Type {
	case "boolean":
		return "false"
	case "integer":
		return "0"
	case "string":
		if !prop.IsRequired {
			return fmt.Sprintf("%q", prop.Name)
		}
		return `""`
	default:
		return `""`
	}
}

func defaultValue(prop *Property) string {
	switch prop.Type {
	case "boolean":
		return "false"
	case "integer":
		return "0"
	case "string":
		if !prop.IsRequired {
			return fmt.Sprintf("stringPtr(%q)", prop.Name)
		}
		return `""`
	default:
		return `""`
	}
}

func downcaseFirst(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	return strings.ToLower(s[0:1]) + s[1:]
}

func getGoType(prop *Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		if prop.IsStruct {
			return "*" + parts[len(parts)-1]
		}
		return parts[len(parts)-1]
	}

	var asterisk string
	if !prop.IsRequired {
		asterisk = "*"
	}

	switch prop.Type {
	case "boolean":
		return asterisk + "bool"
	case "integer":
		return asterisk + "int"
	default:
		return asterisk + prop.Type
	}
}

func multilineComment(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n// ")
}

func optionalMultilineComment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "" // Don't render comment at all
	}
	return "// " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n  // ") + "\n  "
}

func requiredJSONValue(prop *Property, ct *CustomType) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		if prop.IsStruct {
			return "&" + parts[len(parts)-1] + "{}"
		}
		return fmt.Sprintf("%q", prop.FirstEnumValue)
	}

	switch prop.Type {
	case "boolean":
		return "true"
	case "integer":
		return "0"
	case "string":
		return fmt.Sprintf("%q", prop.Name)
	default:
		return `""`
	}
}

func requiredValue(prop *Property, ct *CustomType) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		if prop.IsStruct {
			return "&" + refName + "{}"
		}
		return fmt.Sprintf("%vEnum%v", uppercaseFirst(refName), uppercaseFirst(prop.FirstEnumValue))
	}

	switch prop.Type {
	case "boolean":
		return "true"
	case "integer":
		return "0"
	case "string":
		return fmt.Sprintf("%q", prop.Name)
	default:
		return `""`
	}
}

func showJSONCommaForOptional(index int, ct *CustomType) string {
	if index < len(ct.Properties)-1 {
		return ",\n"
	}
	return ""
}

func showJSONCommaForRequired(index int, ct *CustomType) string {
	for index++; index < len(ct.Properties); index++ {
		prop := ct.Properties[index]
		if prop.IsRequired {
			return ",\n"
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
