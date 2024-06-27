package schema

import (
	"gopkg.in/yaml.v3"
)

func ParseV1(yamlStr string) (*Plugin, error) {
	result := &Plugin{}
	if err := yaml.Unmarshal([]byte(yamlStr), &result); err != nil {
		return nil, err
	}
	return result, nil
}
