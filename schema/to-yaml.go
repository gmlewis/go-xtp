package schema

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func (p *Plugin) ToYaml() (string, error) {
	var buf bytes.Buffer
	e := yaml.NewEncoder(&buf)
	e.SetIndent(2)
	if err := e.Encode(p); err != nil {
		return "", err
	}

	return buf.String(), nil
}
