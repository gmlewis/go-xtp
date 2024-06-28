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
	"addOmitIfNeeded":  addOmitIfNeeded,
	"downcaseFirst":    downcaseFirst,
	"getGoType":        getGoType,
	"multilineComment": multilineComment,
	"uppercaseFirst":   uppercaseFirst,
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

func uppercaseFirst(s string) string {
	if len(s) < 2 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}
