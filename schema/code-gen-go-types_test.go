package schema

import (
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/hand-fruit.go
var wantFruitGo string

//go:embed testdata/hand-fruit_test.go
var wantFruitTestGo string

func TestGenGoCustomTypes(t *testing.T) {
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
			wantSrc:  wantFruitGo,
			wantTest: wantFruitTestGo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin, err := ParseStr(tt.yamlStr)
			if err != nil {
				t.Fatal(err)
			}

			gotSrc, gotTest, err := plugin.genGoCustomTypes()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.wantSrc, gotSrc); diff != "" {
				t.Errorf("genGoCustomTypes src mismatch (-want +got):\n%v", diff)
			}

			if diff := cmp.Diff(tt.wantTest, gotTest); diff != "" {
				t.Errorf("genGoCustomTypes test mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
