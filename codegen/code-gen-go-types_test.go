package codegen

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/gmlewis/go-xtp/schema"
	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/fruit.yaml
var fruitYaml string

//go:embed testdata/user.yaml
var userYaml string

//go:embed testdata/hand-fruit-types.go
var wantFruitGo string

//go:embed testdata/hand-fruit-types_test.go.txt
var wantFruitTestGo string

//go:embed testdata/hand-user-types.go
var wantUserGo string

//go:embed testdata/hand-user-types_test.go.txt
var wantUserTestGo string

func stripLeadingLines(s string, n int) string {
	return strings.Join(strings.Split(s, "\n")[2:], "\n")
}

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
			wantSrc:  stripLeadingLines(wantFruitGo, 2),
			wantTest: stripLeadingLines(wantFruitTestGo, 2),
		},
		{
			name:     "user",
			yamlStr:  userYaml,
			wantSrc:  stripLeadingLines(wantUserGo, 2),
			wantTest: stripLeadingLines(wantUserTestGo, 2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin, err := schema.ParseStr(tt.yamlStr)
			if err != nil {
				t.Fatal(err)
			}

			c := &Client{Plugin: plugin}
			gotSrc, gotTest, err := c.genGoCustomTypes()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.wantSrc, gotSrc); diff != "" {
				t.Logf("got src:\n%v", gotSrc)
				t.Errorf("genGoCustomTypes src mismatch (-want +got):\n%v", diff)
			}

			if diff := cmp.Diff(tt.wantTest, gotTest); diff != "" {
				t.Logf("got test:\n%v", gotTest)
				t.Errorf("genGoCustomTypes test mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
