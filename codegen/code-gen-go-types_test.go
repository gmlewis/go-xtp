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

	genFunc := func(c *Client) (GeneratedFiles, error) {
		if err := c.genGoCustomTypes(); err != nil {
			return nil, err
		}
		m := GeneratedFiles{
			c.CustTypesFilename:      c.CustTypes,
			c.CustTypesTestsFilename: c.CustTypesTests,
		}
		return m, nil
	}

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
			genFunc:     genFunc,
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
			genFunc:     genFunc,
		},
	}

	runEmbedFSTest(t, tests)
}
