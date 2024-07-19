package codegen

import (
	"embed"
	_ "embed"
	"testing"
)

//go:embed testdata/fruit/mbt-types/*
var wantFruitMbtTypesFS embed.FS

//go:embed testdata/user/mbt-types/*
var wantUserMbtTypesFS embed.FS

func TestGenMbtCustomTypes(t *testing.T) {
	t.Parallel()

	tests := []*embedFSTest{
		{
			name:    "fruit",
			lang:    "mbt",
			pkgName: "fruit",
			yamlStr: fruitYaml,
			files: []string{
				"fruit.mbt",
				"fruit_bbtest.mbt",
				"moon.pkg.json",
			},
			embedSubdir: "testdata/fruit/mbt-types",
			embedFS:     wantFruitMbtTypesFS,
			genFunc:     func(c *Client) (GeneratedFiles, error) { return c.GenCustomTypes() },
		},
		{
			name:    "user",
			lang:    "mbt",
			pkgName: "user",
			yamlStr: userYaml,
			files: []string{
				"user.mbt",
				"user_bbtest.mbt",
				"moon.pkg.json",
			},
			embedSubdir: "testdata/user/mbt-types",
			embedFS:     wantUserMbtTypesFS,
			genFunc:     func(c *Client) (GeneratedFiles, error) { return c.GenCustomTypes() },
		},
	}

	runEmbedFSTest(t, tests)
}
