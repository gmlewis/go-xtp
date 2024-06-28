package schema

import (
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/hand-fruit-types.mbt
var wantFruitMbt string

//go:embed testdata/hand-fruit-types_test.mbt
var wantFruitTestMbt string

//go:embed testdata/hand-user-types.mbt
var wantUserMbt string

//go:embed testdata/hand-user-types_test.mbt
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
			plugin, err := ParseStr(tt.yamlStr)
			if err != nil {
				t.Fatal(err)
			}

			gotSrc, gotTest, err := plugin.genMbtCustomTypes()
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
