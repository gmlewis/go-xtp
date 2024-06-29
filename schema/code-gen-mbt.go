package schema

import (
	"fmt"
	"log"
	"strings"
)

func defaultMbtJSONValue(prop *Property, ct *CustomType) string {
	if prop.Ref != "" {
		if !prop.IsRequired && prop.RefCustomType != nil {
			// populate all the required fields recursively:
			requiredProps := prop.RefCustomType.GetRequiredProps()
			fields := make([]string, 0, len(requiredProps))
			for _, p2 := range requiredProps {
				fields = append(fields, fmt.Sprintf("%q:%v", p2.Name, defaultMbtJSONValue(p2, prop.RefCustomType)))
			}
			return fmt.Sprintf("{%v}", strings.Join(fields, ","))
		}
		return fmt.Sprintf("%q", prop.FirstEnumValue)
	}

	switch prop.Type {
	case "integer":
		return "0"
	case "string":
		return `""`
	case "number":
		return "0.0"
	case "boolean":
		return "false"
	case "object":
		return "{}"
	case "array":
		return "[]"
	case "buffer":
		return `""`
	default:
		log.Printf("WARNING: unknown property type %q", prop.Type)
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
	case "integer":
		return "0"
	case "string":
		return `""`
	case "number":
		return "0.0"
	case "boolean":
		return "false"
	case "object":
		return "{}"
	case "array":
		return "[]"
	case "buffer":
		return `""`
	default:
		log.Printf("WARNING: unknown property type %q", prop.Type)
		return `""`
	}
}

func getMbtType(prop *Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		if prop.RefCustomType != nil {
			return refName + "?"
		}
		return refName
	}

	var optional string
	if !prop.IsRequired {
		optional = "?"
	}

	switch prop.Type {
	case "integer":
		return "Int" + optional
	case "string":
		return "String" + optional
	case "number":
		return "Double" + optional
	case "boolean":
		return "Bool" + optional
	case "object":
		return "{}" + optional // TODO - what should this be?
	case "array":
		return "[]" + optional // TODO - what should this be?
	case "buffer":
		return "Buffer" + optional // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", prop.Type)
		return prop.Type + optional
	}
}

func optionalMbtMultilineComment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "" // Don't render comment at all
	}
	return "/// " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n  /// ") + "\n  "
}

func optionalMbtValue(prop *Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		if !prop.IsRequired && prop.RefCustomType != nil {
			return fmt.Sprintf("Some(%v::new())", refName)
		}
		if prop.FirstEnumValue != "" {
			return fmt.Sprintf("Some(%v)", uppercaseFirst(prop.FirstEnumValue))
		}
	}

	switch prop.Type {
	case "integer":
		return "Some(0)"
	case "string":
		return fmt.Sprintf("Some(%q)", prop.Name)
	case "number":
		return "Some(0.0)"
	case "boolean":
		return "Some(true)"
	case "object":
		return "Some({})" // TODO - what should this be?
	case "array":
		return "Some([])" // TODO - what should this be?
	case "buffer":
		return "Some(Buffer)" // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", prop.Type)
		return fmt.Sprintf("Some(%q)", prop.Name)
	}
}

func requiredMbtJSONValue(prop *Property, ct *CustomType) string {
	if prop.Ref != "" {
		if prop.RefCustomType != nil {
			// populate all the required fields recursively:
			requiredProps := prop.RefCustomType.GetRequiredProps()
			fields := make([]string, 0, len(requiredProps))
			for _, p2 := range requiredProps {
				// NOTE: This calls `defaultMbtJSONValue` recursively, not _THIS_ function recursively!
				fields = append(fields, fmt.Sprintf("%q:%v", p2.Name, defaultMbtJSONValue(p2, prop.RefCustomType)))
			}
			return fmt.Sprintf("{%v}", strings.Join(fields, ","))
		}
		return fmt.Sprintf("%q", prop.FirstEnumValue)
	}

	switch prop.Type {
	case "integer":
		return "0"
	case "string":
		return fmt.Sprintf("%q", prop.Name)
	case "number":
		return "0.0"
	case "boolean":
		return "true"
	case "object":
		return "{}" // TODO - what should this be?
	case "array":
		return "[]" // TODO - what should this be?
	case "buffer":
		return "Buffer" // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", prop.Type)
		return `""`
	}
}

func requiredMbtValue(prop *Property) string {
	if !prop.IsRequired {
		return "None"
	}

	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		if prop.RefCustomType != nil {
			return "Some(" + refName + "::new())"
		}
		return uppercaseFirst(prop.FirstEnumValue)
	}

	switch prop.Type {
	case "integer":
		return "0"
	case "string":
		return fmt.Sprintf("%q", prop.Name)
	case "number":
		return "0.0"
	case "boolean":
		return "true"
	case "object":
		return "{}" // TODO - what should this be?
	case "array":
		return "[]" // TODO - what should this be?
	case "buffer":
		return "Buffer" // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", prop.Type)
		return `""`
	}
}
