package schema

import (
	"errors"
	"fmt"
)

var (
	ErrNoCodeGeneration = errors.New("code generation not supported for version v0")
)

// GenCustomTypes generates custom types for the plugin.
func (p *Plugin) GenCustomTypes() (string, error) {
	if p.Version == "v0" {
		return "", ErrNoCodeGeneration
	}

	switch p.Lang {
	case "go":
		return p.genGoCustomTypes()
	case "mbt":
		return p.genMbtCustomTypes()
	default:
		return "", fmt.Errorf("unsupported programming language: %q", p.Lang)
	}
}

// GenHostSDK generates Host SDK code to call the extension plugin.
func (p *Plugin) GenHostSDK(customTypes string) (string, error) {
	if p.Version == "v0" {
		return "", ErrNoCodeGeneration
	}

	switch p.Lang {
	case "go":
		return p.genGoHostSDK(customTypes)
	case "mbt":
		return p.genMbtHostSDK(customTypes)
	default:
		return "", fmt.Errorf("unsupported programming language: %q", p.Lang)
	}
}

// GenPluginPDK generates Plugin PDK code to process plugin calls.
func (p *Plugin) GenPluginPDK(customTypes string) (string, error) {
	if p.Version == "v0" {
		return "", ErrNoCodeGeneration
	}

	switch p.Lang {
	case "go":
		return p.genGoPluginPDK(customTypes)
	case "mbt":
		return p.genMbtPluginPDK(customTypes)
	default:
		return "", fmt.Errorf("unsupported programming language: %q", p.Lang)
	}
}
