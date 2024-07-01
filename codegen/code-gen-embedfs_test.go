package codegen

import (
	"embed"
	"path/filepath"
	"sort"
	"testing"

	"github.com/gmlewis/go-xtp/schema"
	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/fruit.yaml
var fruitYaml string

//go:embed testdata/user.yaml
var userYaml string

type embedFSTest struct {
	name        string
	lang        string
	pkgName     string
	yamlStr     string
	files       []string
	embedSubdir string
	embedFS     embed.FS
	genFunc     func(c *Client) (GeneratedFiles, error)
}

func (e *embedFSTest) readFS(t *testing.T) GeneratedFiles {
	t.Helper()
	if e.embedSubdir == "" {
		t.Fatalf("missing embedSubdir: %+v", *e)
	}

	r := GeneratedFiles{}
	for _, name := range e.files {
		b, err := e.embedFS.ReadFile(filepath.Join(e.embedSubdir, name))
		if err != nil {
			contents, _ := e.embedFS.ReadDir(e.embedSubdir)
			t.Fatalf("%v: files=%+v, contents=%+v: %v", e.embedSubdir, e.files, contents, err)
		}
		r[name] = string(b)
	}
	return r
}

func runEmbedFSTest(t *testing.T, tests []*embedFSTest) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin, err := schema.ParseStr(tt.yamlStr)
			if err != nil {
				t.Fatal(err)
			}

			c, err := New(tt.lang, tt.pkgName, plugin, false)
			if err != nil {
				t.Fatal(err)
			}
			got, err := tt.genFunc(c)
			if err != nil {
				t.Fatal(err)
			}

			wantAll := tt.readFS(t)

			if len(wantAll) != len(got) {
				gotFiles := make([]string, 0, len(got))
				for k := range got {
					gotFiles = append(gotFiles, k)
				}
				sort.Strings(gotFiles)
				t.Errorf("%v generated %v files: %+v, wanted %v files", tt.name, len(got), gotFiles, len(wantAll))
			}

			for _, name := range tt.files {
				want := wantAll[name]
				fullName := filepath.Join(tt.embedSubdir, name)
				if want == "" {
					t.Errorf("test file missing! %v:\n%v", fullName, got[name])
					continue
				}

				if got[name] == "" {
					t.Errorf("file not generated! %v:\n%v", fullName, want)
					continue
				}

				if diff := cmp.Diff(want, got[name]); diff != "" {
					t.Logf("got %v:\n%v", fullName, got[name])
					t.Errorf("gen %q mismatch (-want +got):\n%v", fullName, diff)
				}
			}
		})
	}
}
