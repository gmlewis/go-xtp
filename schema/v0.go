package schema

import (
	"gopkg.in/yaml.v3"
)

type schemaV0 struct {
	Version string   `yaml:"version"`
	Exports []string `yaml:,inline"`
}

func ParseV0(yamlStr string) (*Plugin, error) {
	v0 := &schemaV0{}
	if err := yaml.Unmarshal([]byte(yamlStr), &v0); err != nil {
		return nil, err
	}

	exports := make([]*Export, 0, len(v0.Exports))
	for _, export := range v0.Exports {
		exports = append(exports, &Export{Name: export})
	}
	return &Plugin{Version: v0.Version, Exports: exports}, nil
}
