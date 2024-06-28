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
	"lowerSnakeCase":              lowerSnakeCase,
	"multilineComment":            multilineComment,
	"optionalGoMultilineComment":  optionalGoMultilineComment,
	"optionalMbtMultilineComment": optionalMbtMultilineComment,
	"requiredJSONValue":           requiredJSONValue,
	"requiredValue":               requiredValue,
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

func defaultGoJSONValue(prop *Property, ct *CustomType) string {
	if prop.Ref != "" {
		if !prop.IsRequired && prop.RefCustomType != nil {
			// populate all the required fields recursively:
			requiredProps := prop.RefCustomType.GetRequiredProps()
			fields := make([]string, 0, len(requiredProps))
			for _, p2 := range requiredProps {
				fields = append(fields, fmt.Sprintf("%q:%v", p2.Name, defaultGoJSONValue(p2, prop.RefCustomType)))
			}
			return fmt.Sprintf("{%v}", strings.Join(fields, ","))
		}
		return `""`
	}

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

func defaultGoValue(prop *Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		if !prop.IsRequired && prop.RefCustomType != nil {
			return "&" + refName + "{}"
		}
		return `""`
	}

	switch prop.Type {
	case "boolean":
		if !prop.IsRequired {
			return "boolPtr(false)"
		}
		return "false"
	case "integer":
		if !prop.IsRequired {
			return "intPtr(0)"
		}
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

func defaultMbtJSONValue(prop *Property, ct *CustomType) string {
	if prop.Ref != "" {
		if !prop.IsRequired && prop.RefCustomType != nil {
			// populate all the required fields recursively:
			requiredProps := prop.RefCustomType.GetRequiredProps()
			fields := make([]string, 0, len(requiredProps))
			for _, p2 := range requiredProps {
				fields = append(fields, fmt.Sprintf("%q:%v", p2.Name, defaultGoJSONValue(p2, prop.RefCustomType)))
			}
			return fmt.Sprintf("{%v}", strings.Join(fields, ","))
		}
		return `""`
	}

	switch prop.Type {
	case "boolean":
		return "false"
	case "integer":
		return "0"
	case "string":
		return `""`
	default:
		return `""`
	}
}

func defaultMbtValue(prop *Property) string {
	if prop.Ref != "" {
		// parts := strings.Split(prop.Ref, "/")
		// refName := parts[len(parts)-1]
		if !prop.IsRequired && prop.RefCustomType != nil {
			return "None"
		}
		if prop.FirstEnumValue != "" {
			return uppercaseFirst(prop.FirstEnumValue)
		}
		return `""`
	}

	if !prop.IsRequired {
		return "None"
	}

	switch prop.Type {
	case "boolean":
		return "false"
	case "integer":
		return "0"
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
		if prop.RefCustomType != nil {
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

func getMbtType(prop *Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		if prop.RefCustomType != nil {
			return parts[len(parts)-1] + "?"
		}
		return parts[len(parts)-1]
	}

	var optional string
	if !prop.IsRequired {
		optional = "?"
	}

	switch prop.Type {
	case "boolean":
		return "Bool" + optional
	case "integer":
		return "Int" + optional
	case "string":
		return "String" + optional
	default:
		return prop.Type + optional
	}
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

func optionalGoMultilineComment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "" // Don't render comment at all
	}
	return "// " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n  // ") + "\n  "
}

func optionalMbtMultilineComment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "" // Don't render comment at all
	}
	return "/// " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n  /// ") + "\n  "
}

func requiredJSONValue(prop *Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		if prop.RefCustomType != nil {
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

func requiredValue(prop *Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		if prop.RefCustomType != nil {
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
