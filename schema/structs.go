package schema

import "errors"

var (
	ErrNoCodeGeneration = errors.New("code generation not supported for version v0")
)

// Structs represents Go structs that access the XTP Extension Plugin functions.
type Structs struct {
}

func (s *Structs) String() string {
	return ""
}

func (p *Plugin) GenStructs() (*Structs, error) {
	if p.Version == "v0" {
		return nil, ErrNoCodeGeneration
	}
	return nil, nil
}

func (p *Plugin) GenMain(structs *Structs) (string, error) {
	if p.Version == "v0" {
		return "", ErrNoCodeGeneration
	}
	return "", nil
}
