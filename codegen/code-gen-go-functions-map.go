package codegen

import (
	"fmt"
	"log"
	"strings"

	"github.com/gmlewis/go-xtp/schema"
)

func defaultGoJSONValue(prop *schema.Property, ct *schema.CustomType) string {
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
	case "integer":
		return "0"
	case "string":
		if !prop.IsRequired {
			return fmt.Sprintf("%q", prop.Name)
		}
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

func defaultGoValue(prop *schema.Property) string {
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

func getGoType(prop *schema.Property) string {
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
	case "integer":
		return asterisk + "int"
	case "string":
		return asterisk + "string"
	case "number":
		return asterisk + "float64"
	case "boolean":
		return asterisk + "bool"
	case "object":
		return asterisk + "{}" // TODO - what should this be?
	case "array":
		return asterisk + "[]" // TODO - what should this be?
	case "buffer":
		return asterisk + "buffer" // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", prop.Type)
		return asterisk + prop.Type
	}
}

func goMultilineComment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "\n"
	}
	return "// " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n  // ")
}

func inputToGoType(input *schema.Input) string {
	if input == nil {
		return ""
	}

	if input.Ref != "" {
		parts := strings.Split(input.Ref, "/")
		refName := parts[len(parts)-1]
		return "input " + refName
	}

	switch input.Type {
	case "integer":
		return "input int"
	case "string":
		return "input string"
	case "number":
		return "input float64"
	case "boolean":
		return "input bool"
	case "object":
		return "input {}" // TODO - what should this be?
	case "array":
		return "input []" // TODO - what should this be?
	case "buffer":
		return "input Buffer" // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", input.Type)
		return "input " + input.Type
	}
}

func jsonOutputAsGoType(output *schema.Output) string {
	// TODO: finish this
	// if output.Ref != "" {
	// 	parts := strings.Split(output.Ref, "/")
	// 	refName := parts[len(parts)-1]
	// 	if output.RefCustomType != nil {
	// 		return refName + "?"
	// 	}
	// 	return refName
	// }

	// var optional string
	// if !output.IsRequired {
	// 	optional = "?"
	// }

	switch output.Type {
	// case "integer":
	// 	return "Int" + optional
	// case "string":
	// 	return "String" + optional
	// case "number":
	// 	return "Double" + optional
	case "boolean":
		return "bool"
	// case "object":
	// 	return "{}" + optional // TODO - what should this be?
	// case "array":
	// 	return "[]" + optional // TODO - what should this be?
	// case "buffer":
	// 	return "Buffer" + optional // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", output.Type)
		return output.Type
	}
}

func optionalGoMultilineComment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "" // Don't render comment at all
	}
	return "// " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n  // ") + "\n  "
}

func outputToGoExampleLiteral(output *schema.Output) string {
	if output == nil {
		return ""
	}

	if output.Ref != "" {
		parts := strings.Split(output.Ref, "/")
		refName := parts[len(parts)-1]
		return fmt.Sprintf("\n\treturn %v{}", refName)
	}

	switch output.Type {
	case "integer":
		return "\n\treturn 0"
	case "string":
		return `\n\treturn ""`
	case "number":
		return "\n\treturn 0.0"
	case "boolean":
		return "\n\treturn false"
	case "object":
		return "\n\treturn {}" // TODO - what should this be?
	case "array":
		return "\n\treturn []" // TODO - what should this be?
	case "buffer":
		return "\n\treturn Buffer" // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", output.Type)
		return "\n\t" + output.Type
	}
}

func outputToGoType(output *schema.Output) string {
	if output == nil {
		return ""
	}

	if output.Ref != "" {
		parts := strings.Split(output.Ref, "/")
		refName := parts[len(parts)-1]
		// if output.RefCustomType != nil {
		// 	return refName + "?"
		// }
		return refName
	}

	switch output.Type {
	case "integer":
		return "int"
	case "string":
		return "string"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "object":
		return "{}" // TODO - what should this be?
	case "array":
		return "[]" // TODO - what should this be?
	case "buffer":
		return "Buffer" // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", output.Type)
		return output.Type
	}
}

func requiredGoJSONValue(prop *schema.Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		if prop.RefCustomType != nil {
			return "&" + parts[len(parts)-1] + "{}" // FIX - not a JSON value
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

func requiredGoValue(prop *schema.Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		if prop.RefCustomType != nil {
			return fmt.Sprintf("&%v{}", refName)
		}
		return fmt.Sprintf("%vEnum%v", uppercaseFirst(refName), uppercaseFirst(prop.FirstEnumValue))
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
