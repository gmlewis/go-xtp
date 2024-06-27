package schema

import "errors"

var (
	ErrNoCodeGeneration = errors.New("code generation not supported for version v0")
)

func (p *Plugin) GenStructs() (string, error) {
	if p.Version == "v0" {
		return "", ErrNoCodeGeneration
	}
	return "", nil
}

func (p *Plugin) GenMain(structs string) (string, error) {
	if p.Version == "v0" {
		return "", ErrNoCodeGeneration
	}
	return "", nil
}
