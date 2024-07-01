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

	genFunc := func(c *Client) (GeneratedFiles, error) {
		if err := c.genMbtCustomTypes(); err != nil {
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
			lang:    "mbt",
			pkgName: "fruit",
			yamlStr: fruitYaml,
			files: []string{
				"fruit.mbt",
				"fruit_test.mbt",
			},
			embedSubdir: "testdata/fruit/mbt-types",
			embedFS:     wantFruitMbtTypesFS,
			genFunc:     genFunc,
		},
		{
			name:    "user",
			lang:    "mbt",
			pkgName: "user",
			yamlStr: userYaml,
			files: []string{
				"user.mbt",
				"user_test.mbt",
			},
			embedSubdir: "testdata/user/mbt-types",
			embedFS:     wantUserMbtTypesFS,
			genFunc:     genFunc,
		},
	}

	runEmbedFSTest(t, tests)
}
