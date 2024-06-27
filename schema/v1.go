package schema

import (
	"gopkg.in/yaml.v3"
)

func ParseV1(yamlStr string) (*Plugin, error) {
	result := &Plugin{}
	if err := yaml.Unmarshal([]byte(yamlStr), &result); err != nil {
		return nil, err
	}

	// run through all the Properties and populate the "IsRequired" field as needed.
	for _, ct := range result.CustomTypes {
		reqFields := map[string]bool{}
		for _, req := range ct.Required {
			reqFields[req] = true
		}
		for _, prop := range ct.Properties {
			if reqFields[prop.Name] {
				prop.IsRequired = true
			}
		}
	}

	return result, nil
}
