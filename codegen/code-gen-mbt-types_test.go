package codegen

import (
	_ "embed"
	"testing"

	"github.com/gmlewis/go-xtp/schema"
	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/fruit/mbt-types/fruit.mbt
var wantFruitMbt string

//go:embed testdata/fruit/mbt-types/fruit_test.mbt
var wantFruitTestMbt string

//go:embed testdata/user/mbt-types/user.mbt
var wantUserMbt string

//go:embed testdata/user/mbt-types/user_test.mbt
var wantUserTestMbt string

func TestGenMbtCustomTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		yamlStr  string
		wantSrc  string
		wantTest string
	}{
		{
			name:     "fruit",
			yamlStr:  fruitYaml,
			wantSrc:  wantFruitMbt,
			wantTest: wantFruitTestMbt,
		},
		{
			name:     "user",
			yamlStr:  userYaml,
			wantSrc:  wantUserMbt,
			wantTest: wantUserTestMbt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin, err := schema.ParseStr(tt.yamlStr)
			if err != nil {
				t.Fatal(err)
			}

			c := &Client{Plugin: plugin}
			gotSrc, gotTest, err := c.genMbtCustomTypes()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.wantSrc, gotSrc); diff != "" {
				t.Logf("got src:\n%v", gotSrc)
				t.Errorf("genMbtCustomTypes src mismatch (-want +got):\n%v", diff)
			}

			if diff := cmp.Diff(tt.wantTest, gotTest); diff != "" {
				t.Logf("got test:\n%v", gotTest)
				t.Errorf("genMbtCustomTypes test mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
