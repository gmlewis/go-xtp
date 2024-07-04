package codegen

import (
	"embed"
	_ "embed"
	"testing"
)

//go:embed testdata/fruit/mbt-plugin/*
var wantFruitMbtPluginFS embed.FS

//go:embed testdata/user/mbt-plugin/*
var wantUserMbtPluginFS embed.FS

func TestGenMbtPluginPDK(t *testing.T) {
	t.Parallel()
	tests := []*embedFSTest{
		{
			name:    "fruit",
			lang:    "mbt",
			pkgName: "fruit",
			yamlStr: fruitYaml,
			files: []string{
				"build.sh",
				"fruit.mbt",
				"fruit_test.mbt",
				"host-functions.mbt",
				"main.mbt",
				"moon.pkg.json",
				"plugin-functions.mbt",
				"xtp.toml",
			},
			embedSubdir: "testdata/fruit/mbt-plugin",
			embedFS:     wantFruitMbtPluginFS,
			genFunc:     func(c *Client) (GeneratedFiles, error) { return c.genMbtPluginPDK() },
		},
		{
			name:    "user",
			lang:    "mbt",
			pkgName: "user",
			yamlStr: userYaml,
			files: []string{
				"build.sh",
				"user.mbt",
				"user_test.mbt",
				"main.mbt",
				"moon.pkg.json",
				"plugin-functions.mbt",
				"xtp.toml",
			},
			embedSubdir: "testdata/user/mbt-plugin",
			embedFS:     wantUserMbtPluginFS,
			genFunc:     func(c *Client) (GeneratedFiles, error) { return c.genMbtPluginPDK() },
		},
	}

	runEmbedFSTest(t, tests)
}
