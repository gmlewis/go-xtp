package schema

import (
	_ "embed"
	"testing"
)

//go:embed testdata/hand-fruit.mbt
var wantFruitMbt string

//go:embed testdata/hand-fruit_test.mbt
var wantFruitTestMbt string

func TestGenMbtCustomTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
	}{
		{
			name: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
