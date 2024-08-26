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

func getMbtType(item any) string {
	var ref string
	var isRequired bool
	var itemType string
	var refCustomType *schema.CustomType

	switch t := item.(type) {
	case *schema.Property:
		ref = t.Ref
		isRequired = t.IsRequired
		itemType = t.Type
		refCustomType = t.RefCustomType
	case *schema.Output:
		ref = t.Ref
		itemType = t.Type
		isRequired = true
	default:
		log.Fatalf("getMbtType: unsupported type: %T", t)
	}

	if ref != "" {
		parts := strings.Split(ref, "/")
		refName := parts[len(parts)-1]
		if refCustomType != nil {
			return refName + "?"
		}
		return refName
	}

	var optional string
	if !isRequired {
		optional = "?"
	}

	switch itemType {
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
		log.Printf("WARNING: unknown property type %q", itemType)
		return itemType + optional
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

func mbtTypeIs(item any, name string) bool {
	mbtType := getMbtType(item)
	return mbtType == name
}

func mbtTypeIsOptional(prop *schema.Property) bool {
	mbtType := getMbtType(prop)
	return strings.HasSuffix(mbtType, "?")
}

// This function's output values matches the output from optionalMbtValue.
func optionalMbtJSONValue(prop *schema.Property, ct *schema.CustomType) string {
	if prop.Ref != "" {
		if prop.RefCustomType != nil {
			// populate all the required fields recursively:
			requiredProps := prop.RefCustomType.GetRequiredProps()
			fields := make([]string, 0, len(requiredProps))
			for _, p2 := range requiredProps {
				fields = append(fields, fmt.Sprintf("%q:%v", p2.Name, optionalMbtJSONValue(p2, prop.RefCustomType)))
			}
			return fmt.Sprintf("{%v}", strings.Join(fields, ","))
		}
		return fmt.Sprintf("%q", prop.FirstEnumValue)
	}

	return mbtTypeTestValue(prop.Type, prop.Name, prop.IsRequired)
}

func optionalMbtMultilineComment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "" // Don't render comment at all
	}
	return "/// " + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n  /// ") + "\n  "
}

// If a property is required, use that type's default value, otherwise return a non-default value.
func mbtTypeTestValue(propType, propName string, isRequired bool) string {
	switch propType {
	case "integer":
		if isRequired {
			return "0"
		}
		return "42"
	case "string":
		if isRequired {
			return `""`
		}
		return fmt.Sprintf("%q", propName)
	case "number":
		if isRequired {
			return "0.0"
		}
		return "42.0"
	case "boolean":
		if isRequired {
			return "false"
		}
		return "true"
	case "object":
		return "{}" // TODO - what should this be?
	case "array":
		return "[]" // TODO - what should this be?
	case "buffer":
		return "Buffer" // TODO - what should this be?
	default:
		log.Printf("WARNING: unknown property type %q", propType)
		return fmt.Sprintf("%q", propName)
	}
}

// This function's output values matches the output from optionalMbtJSONValue.
func optionalMbtValue(prop *schema.Property, ct *schema.CustomType) string {
	if prop.Ref != "" {
		if prop.RefCustomType != nil {
			// populate all the required fields recursively:
			requiredProps := prop.RefCustomType.GetRequiredProps()
			fields := make([]string, 0, len(requiredProps))
			for _, p2 := range requiredProps {
				fields = append(fields, fmt.Sprintf("%v: %v", p2.Name, optionalMbtValue(p2, prop.RefCustomType)))
			}
			if !prop.IsRequired {
				return fmt.Sprintf("Some({%v})", strings.Join(fields, ","))
			}
			return fmt.Sprintf("{%v}", strings.Join(fields, ","))
		}
		return fmt.Sprintf("%q", prop.FirstEnumValue)
	}

	value := mbtTypeTestValue(prop.Type, prop.Name, prop.IsRequired)
	if prop.IsRequired {
		return value
	}
	return fmt.Sprintf("Some(%v)", value)
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
