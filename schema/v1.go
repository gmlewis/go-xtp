package schema

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseV1(yamlStr string) (*Plugin, error) {
	result := &Plugin{}
	if err := yaml.Unmarshal([]byte(yamlStr), &result); err != nil {
		return nil, err
	}

	// run through all the Properties and populate the "IsRequired"
	// and "IsStruct" fields as needed.
	isStruct := map[string]bool{}
	firstEnumValue := map[string]string{}
	for _, ct := range result.CustomTypes {
		if len(ct.Enum) == 0 && len(ct.Properties) > 0 {
			isStruct[ct.Name] = true
		}
		if len(ct.Enum) > 0 {
			firstEnumValue[ct.Name] = ct.Enum[0]
		}
	}

	for _, ct := range result.CustomTypes {
		reqFields := map[string]bool{}
		for _, req := range ct.Required {
			reqFields[req] = true
		}
		for _, prop := range ct.Properties {
			if reqFields[prop.Name] {
				prop.IsRequired = true
			}
			if prop.Ref != "" {
				parts := strings.Split(prop.Ref, "/")
				refName := parts[len(parts)-1]
				if isStruct[refName] {
					prop.IsStruct = true
				}
				if v := firstEnumValue[refName]; v != "" {
					prop.FirstEnumValue = v
				}
			}
		}
	}

	return result, nil
}
