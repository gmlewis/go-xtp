package codegen

import (
	"fmt"
	"log"
	"strings"

	"github.com/gmlewis/go-xtp/schema"
)

func defaultMbtJSONValue(prop *schema.Property, ct *schema.CustomType) string {
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

func defaultMbtValue(prop *schema.Property) string {
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

func getMbtType(prop *schema.Property) string {
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

func inputToMbtType(input *schema.Input) string {
	if input == nil {
		return ""
	}

	if input.Ref != "" {
		parts := strings.Split(input.Ref, "/")
		refName := parts[len(parts)-1]
		return "input : " + refName
	}

	switch input.Type {
	case "integer":
		return "input : Int"
	case "string":
		return "input : String"
	case "number":
		return "input : Double"
	case "boolean":
		return "input : Bool"
	case "object":
		return "input : {}" // TODO - what should this be?
	case "array":
		return "input : []" // TODO - what should this be?
	case "buffer":
		return "input : Buffer" // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", input.Type)
		return "input : " + input.Type
	}
}

func jsonOutputAsMbtType(output *schema.Output) string {
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
		return "as_bool"
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

func mbtConvertFromJSONValue(prop *schema.Property) string {
	valueGet := fmt.Sprintf("value.get(%q)", prop.Name)
	if prop.IsRequired {
		valueGet += "?" // fail shortcut if not Some()
	}

	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		if !prop.IsRequired {
			return fmt.Sprintf(`match %v {
    Some(jv) => %v::from_json(jv)
    None => None
  }`, valueGet, refName)
		}
		return fmt.Sprintf("%v |> %v::from_json()", valueGet, refName)
	}

	var asType string
	switch prop.Type {
	case "integer":
		asType = ".as_number()"
	case "string":
		asType = ".as_string()"
	case "number":
		asType = ".as_number()"
	case "boolean":
		asType = ".as_bool()"
	default:
		log.Printf("WARNING: unknown property type %q", prop.Type)
		return `"unknown"`
	}

	if !prop.IsRequired {
		switch prop.Type {
		case "integer":
			return fmt.Sprintf(`match %v {
    Some(jv) => json_as_integer(jv)
    None => None
  }`, valueGet)
		default:
			return fmt.Sprintf(`match %v {
    Some(jv) => jv%v
    None => None
  }`, valueGet, asType)
		}
	}

	switch prop.Type {
	case "integer":
		return fmt.Sprintf("json_as_integer(%v)", valueGet)
	default:
		return valueGet + asType
	}
}

func mbtFromJSONMatchKey(prop *schema.Property) string {
	if prop.Ref != "" {
		// parts := strings.Split(prop.Ref, "/")
		// refName := parts[len(parts)-1]
		// if !prop.IsRequired && prop.RefCustomType != nil {
		// 	return "None"
		// }
		// if prop.FirstEnumValue != "" {
		// 	return uppercaseFirst(prop.FirstEnumValue)
		// }
		return "v"
	}

	// if !prop.IsRequired {
	// 	return "None"
	// }

	switch prop.Type {
	case "integer", "number":
		return "Json::Number(n)"
	case "string":
		return "Json::String(s)"
	case "boolean":
		return "Json::Boolean(v)"
	case "object":
		return "{}" // TODO - what to do with this?
	case "array":
		return "[]" // TODO - what to do with this?
	case "buffer":
		return `""` // TODO - what to do with this?
	default:
		log.Printf("WARNING: unknown property type %q", prop.Type)
		return `""`
	}
}

func mbtFromJSONMatchValue(prop *schema.Property) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
		refName := parts[len(parts)-1]
		// if !prop.IsRequired && prop.RefCustomType != nil {
		// 	return "None"
		// }
		// if prop.FirstEnumValue != "" {
		// 	return uppercaseFirst(prop.FirstEnumValue)
		// }
		return fmt.Sprintf("%v::from_json(v)", refName)
	}

	// if !prop.IsRequired {
	// 	return "None"
	// }

	switch prop.Type {
	case "integer":
		return "Some(n.to_int())"
	case "string":
		return "Some(s.to_string())"
	case "number":
		return "Some(n.to_double())"
	case "boolean":
		return "Some(v.to_bool())"
	case "object":
		return "{}" // TODO - what to do with this?
	case "array":
		return "[]" // TODO - what to do with this?
	case "buffer":
		return `""` // TODO - what to do with this?
	default:
		log.Printf("WARNING: unknown property type %q", prop.Type)
		return `""`
	}
}

func mbtMultilineComment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "\n"
	}
	return "/// " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n  /// ")
}

func optionalMbtMultilineComment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "" // Don't render comment at all
	}
	return "/// " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n  /// ") + "\n  "
}

func optionalMbtValue(prop *schema.Property) string {
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

func outputToMbtExampleLiteral(output *schema.Output) string {
	if output == nil {
		return ""
	}

	if output.Ref != "" {
		parts := strings.Split(output.Ref, "/")
		refName := parts[len(parts)-1]
		return fmt.Sprintf(`
  {
    ..%v::new(),
  }`, refName)
	}

	switch output.Type {
	case "integer":
		return "\n  0"
	case "string":
		return `\n  ""`
	case "number":
		return "\n  0.0"
	case "boolean":
		return "\n  false"
	case "object":
		return "\n  {}" // TODO - what should this be?
	case "array":
		return "\n  []" // TODO - what should this be?
	case "buffer":
		return "\n  Buffer" // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", output.Type)
		return "\n  " + output.Type
	}
}

func outputToMbtType(output *schema.Output) string {
	if output == nil {
		return "Unit"
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
		return "Int"
	case "string":
		return "String"
	case "number":
		return "Double"
	case "boolean":
		return "Bool"
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

func requiredMbtJSONValue(prop *schema.Property, ct *schema.CustomType) string {
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

func requiredMbtValue(prop *schema.Property) string {
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
