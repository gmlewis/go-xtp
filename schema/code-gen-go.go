package schema

import (
	"fmt"
	"strings"
)

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
			return fmt.Sprintf("&%v{}", refName)
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

func getGoType(prop *Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		if prop.RefCustomType != nil {
			return "*" + refName
		}
		return refName
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

func optionalGoMultilineComment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "" // Don't render comment at all
	}
	return "// " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n  // ") + "\n  "
}

func requiredGoJSONValue(prop *Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		if prop.RefCustomType != nil {
			return "&" + parts[len(parts)-1] + "{}" // FIX - not a JSON value
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

func requiredGoValue(prop *Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		if prop.RefCustomType != nil {
			return fmt.Sprintf("&%v{}", refName)
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
