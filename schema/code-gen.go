package schema

import (
	"errors"
	"fmt"
)

var (
	ErrNoCodeGeneration = errors.New("code generation not supported for version v0")
)

// GenCustomTypes generates custom types with tests for the plugin.
func (p *Plugin) GenCustomTypes() (srcFile, testFile string, err error) {
	if p.Version == "v0" {
		return "", "", ErrNoCodeGeneration
	}

	switch p.Lang {
	case "go":
		return p.genGoCustomTypes()
	case "mbt":
		return p.genMbtCustomTypes()
	default:
		return "", "", fmt.Errorf("unsupported programming language: %q", p.Lang)
	}
}

// GenHostSDK generates Host SDK code to call the extension plugin.
func (p *Plugin) GenHostSDK(customTypes, pkgName string) (string, error) {
	if p.Version == "v0" {
		return "", ErrNoCodeGeneration
	}

	switch p.Lang {
	case "go":
		return p.genGoHostSDK(customTypes, pkgName)
	case "mbt":
		return p.genMbtHostSDK(customTypes, pkgName)
	default:
		return "", fmt.Errorf("unsupported programming language: %q", p.Lang)
	}
}

// GenPluginPDK generates Plugin PDK code to process plugin calls.
func (p *Plugin) GenPluginPDK(customTypes, pkgName string) (string, error) {
	if p.Version == "v0" {
		return "", ErrNoCodeGeneration
	}

	switch p.Lang {
	case "go":
		return p.genGoPluginPDK(customTypes, pkgName)
	case "mbt":
		return p.genMbtPluginPDK(customTypes, pkgName)
	default:
		return "", fmt.Errorf("unsupported programming language: %q", p.Lang)
	}
}
