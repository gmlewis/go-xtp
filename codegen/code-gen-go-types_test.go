package codegen

import (
	"embed"
	"testing"
)

//go:embed testdata/fruit/go-types/*
var wantFruitGoTypesFS embed.FS

//go:embed testdata/user/go-types/*
var wantUserGoTypesFS embed.FS

func TestGenGoCustomTypes(t *testing.T) {
	t.Parallel()

	tests := []*embedFSTest{
		{
			name:    "fruit",
			lang:    "go",
			pkgName: "fruit",
			yamlStr: fruitYaml,
			files: []string{
				"fruit.go",
				"fruit_test.go",
			},
			embedSubdir: "testdata/fruit/go-types",
			embedFS:     wantFruitGoTypesFS,
			genFunc:     func(c *Client) (GeneratedFiles, error) { return c.GenCustomTypes() },
		},
		{
			name:    "user",
			lang:    "go",
			pkgName: "user",
			yamlStr: userYaml,
			files: []string{
				"user.go",
				"user_test.go",
			},
			embedSubdir: "testdata/user/go-types",
			embedFS:     wantUserGoTypesFS,
			genFunc:     func(c *Client) (GeneratedFiles, error) { return c.GenCustomTypes() },
		},
	}

	runEmbedFSTest(t, tests)
}
