package schema

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/hand-fruit.go
var wantFruitGo string

//go:embed testdata/hand-fruit_test.go
var wantFruitTestGo string

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
				t.Log(gotSrc)
				t.Errorf("genGoCustomTypes src mismatch (-want +got):\n%v", diff)
			}

			if diff := cmp.Diff(tt.wantTest, gotTest); diff != "" {
				t.Log(gotTest)
				t.Errorf("genGoCustomTypes test mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestGenGoCustomType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ct   *CustomType
		want string
	}{
		{
			name: "enum",
			ct: &CustomType{
				Name:        "Fruit",
				Description: "A set of available fruits you can consume",
				Enum:        []string{"apple", "orange", "banana", "strawberry"},
			},
			want: `// Fruit represents a set of available fruits you can consume.
type Fruit string

const (
	FruitEnumApple      Fruit = "apple"
	FruitEnumOrange     Fruit = "orange"
	FruitEnumBanana     Fruit = "banana"
	FruitEnumStrawberry Fruit = "strawberry"
)
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := &Plugin{CustomTypes: []*CustomType{tt.ct}}

			got, err := plugin.genGoCustomType(tt.ct)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Log(got)
				t.Errorf("genGoCustomType mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
