package pet

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k1LoW/frgm/snippet"
)

func TestLoadFile(t *testing.T) {
	tests := []struct {
		in   string
		want Snippets
	}{
		{
			in: `
`,
			want: Snippets{},
		},
		{
			in: `[[snippets]]
  description = "ping"
  command = "ping 8.8.8.8"
  tag = ["network", "google"]
  output = ""

[[snippets]]
  description = "hello"
  command = "echo hello"
  tag = ["sample"]
  output = "hello"
`,
			want: Snippets{
				Snippets: []Snippet{
					Snippet{
						Command:     "ping 8.8.8.8",
						Description: "ping",
						Output:      "",
						Tag:         []string{"network", "google"},
					},
					Snippet{
						Command:     "echo hello",
						Description: "hello",
						Output:      "hello",
						Tag:         []string{"sample"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		in := bytes.NewBufferString(tt.in)
		pet := New()
		got, err := pet.LoadFile(in)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Errorf("got %v\nwant %v", got, tt.want)
		}
	}
}

func TestExport(t *testing.T) {
	uid, err := genUID("ping 8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		current string
		in      snippet.Snippets
		want    string
	}{
		{
			current: `[[snippets]]
  description = "ping"
  command = "ping 8.8.8.8"
  tag = ["network", "google"]
  output = ""
`,
			in: snippet.Snippets{
				snippet.New("uid", "group", "name", "content", "output", "description", []string{"a", "b"}),
			},
			want: `[[snippets]]
  command = "ping 8.8.8.8"
  description = "ping"
  output = ""
  tag = ["network", "google"]

[[snippets]]
  command = "content"
  description = "name"
  output = "output"
  tag = ["a", "b"]
`,
		},
		{
			current: `[[snippets]]
  description = "ping"
  command = "ping 8.8.8.8"
  tag = ["network", "google"]
  output = ""
`,
			in: snippet.Snippets{
				snippet.New(uid, "group", "name", "content", "output", "description", []string{"a", "b"}),
			},
			want: `[[snippets]]
  command = "ping 8.8.8.8"
  description = "name"
  output = "output"
  tag = ["a", "b"]
`,
		},
	}

	tmp, err := ioutil.TempDir("", "frgm")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	for _, tt := range tests {
		pet := New()
		dest := filepath.Join(tmp, "pet.toml")
		if err := ioutil.WriteFile(dest, []byte(tt.current), 0644); err != nil {
			t.Fatal(err)
		}
		if err := pet.Export(tt.in, dest); err != nil {
			t.Fatal(err)
		}
		d, err := ioutil.ReadFile(dest)
		if err != nil {
			t.Fatal(err)
		}
		got := string(d)
		if got != tt.want {
			t.Errorf("got %v\nwant %v", got, tt.want)
		}
	}
}
