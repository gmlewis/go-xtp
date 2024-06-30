package testdata // this line will be stripped in the unit test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

func boolPtr(b bool) *bool       { return &b }
func intPtr(i int) *int          { return &i }
func stringPtr(s string) *string { return &s }

func TestAddressMarshal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		obj  *Address
		want string
	}{
		{
			name: "required fields",
			obj: &Address{
				Street: "street",
			},
			want: `{"street":"street"}`,
		},
		{
			name: "optional fields",
			obj:  &Address{},
			want: `{"street":""}`,
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

func TestUserMarshal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		obj  *User
		want string
	}{
		{
			name: "required fields",
			obj:  &User{},
			want: `{}`,
		},
		{
			name: "optional fields",
			obj: &User{
				Age:     intPtr(0),
				Email:   stringPtr("email"),
				Address: &Address{},
			},
			want: `{"age":0,"email":"email","address":{"street":""}}`,
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
