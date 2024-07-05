package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
