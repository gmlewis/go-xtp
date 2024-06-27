package schema

import "strings"

// genGoCustomTypes generates custom types with tests for the plugin in Go.
func (p *Plugin) genGoCustomTypes() (srcFile, testFile string, err error) {
	var srcBlocks, testBlocks []string
	return strings.Join(srcBlocks, "\n\n"), strings.Join(testBlocks, "\n\n"), nil
}
