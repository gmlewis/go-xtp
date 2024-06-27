package fruit // this line will be stripped in the unit test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

func stringPtr(s string) *string { return &s }

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
				t.Log(string(got))
				t.Errorf("Marshal mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
