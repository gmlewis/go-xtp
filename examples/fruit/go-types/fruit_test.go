package fruit

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

func boolPtr(b bool) *bool       { return &b }
func intPtr(i int) *int          { return &i }
func stringPtr(s string) *string { return &s }

func TestParseFruit(t *testing.T) {
	t.Parallel()

	fruit := FruitEnumApple
	buf, err := jsoncomp.Marshal(fruit)
	if err != nil {
		t.Fatal(err)
	}

	want := `"apple"`
	if got := string(buf); got != want {
		t.Errorf("Marshal = '%v', want '%v'", got, want)
	}

	got, err := ParseFruit(want)
	if err != nil {
		t.Fatal(err)
	}
	if got != fruit {
		t.Errorf("ParseFruit = '%v', want '%v'", got, fruit)
	}
}

func TestParseGhostGang(t *testing.T) {
	t.Parallel()

	ghostGang := GhostGangEnumBlinky
	buf, err := jsoncomp.Marshal(ghostGang)
	if err != nil {
		t.Fatal(err)
	}

	want := `"blinky"`
	if got := string(buf); got != want {
		t.Errorf("Marshal = '%v', want '%v'", got, want)
	}

	got, err := ParseGhostGang(want)
	if err != nil {
		t.Fatal(err)
	}
	if got != ghostGang {
		t.Errorf("ParseGhostGang = '%v', want '%v'", got, ghostGang)
	}
}

func TestComplexObjectMarshal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		obj  *ComplexObject
		want string
	}{
		{
			name: "required fields",
			obj: &ComplexObject{
				Ghost:    GhostGangEnumBlinky,
				ABoolean: true,
				AString:  "aString",
				AnInt:    0,
			},
			want: `{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0}`,
		},
		{
			name: "optional fields",
			obj: &ComplexObject{
				AnOptionalDate: stringPtr("anOptionalDate"),
			},
			want: `{"ghost":"","aBoolean":false,"aString":"","anInt":0,"anOptionalDate":"anOptionalDate"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsoncomp.Marshal(tt.obj)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.want, string(got)); diff != "" {
				t.Logf("got:\n%v", string(got))
				t.Errorf("Marshal mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
