package codegen

import (
	"embed"
	"testing"
)

//go:embed testdata/fruit/go-plugin/*
var wantFruitGoPluginFS embed.FS

//go:embed testdata/user/go-plugin/*
var wantUserGoPluginFS embed.FS

func TestGenGoPluginPDK(t *testing.T) {
	t.Parallel()
	tests := []*embedFSTest{
		{
			name:    "fruit",
			yamlStr: fruitYaml,
			files: []string{
				"build.sh",
				"fruit.go",
				"fruit_test.go",
				"host-functions.go",
				"main.go",
				"plugin-functions.go",
				"xtp.toml",
			},
			embedSubdir: "testdata/fruit/go-plugin",
			embedFS:     wantFruitGoPluginFS,
			genFunc:     func(c *Client) (GeneratedFiles, error) { return c.genGoPluginPDK() },
		},
		// {
		// 	name:     "user",
		// 	yamlStr:  userYaml,
		// 	wantSrc:  wantUserGo,
		// 	wantTest: wantUserTestGo
		// },
	}

	runEmbedFSTest(t, tests)
}
