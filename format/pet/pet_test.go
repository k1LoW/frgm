package pet

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k1LoW/frgm/snippet"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		in   string
		want snippet.Snippets
	}{
		{
			in: `
`,
			want: snippet.Snippets{},
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
			want: snippet.Snippets{
				&snippet.Snippet{
					UID:     "pet-df7bb29f9681",
					Group:   "default",
					Content: "ping 8.8.8.8",
					Name:    "ping",
					Output:  "",
					Labels:  []string{"network", "google"},
				},
				&snippet.Snippet{
					UID:     "pet-584a331fd6b0",
					Group:   "default",
					Content: "echo hello",
					Name:    "hello",
					Output:  "hello",
					Labels:  []string{"sample"},
				},
			},
		},
	}
	for _, tt := range tests {
		in := bytes.NewBufferString(tt.in)
		pet := New()
		got, err := pet.Decode(in, "default")
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

	tmp, err := os.MkdirTemp("", "frgm")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	for _, tt := range tests {
		pet := New()
		dest := filepath.Join(tmp, "pet.toml")
		if err := os.WriteFile(dest, []byte(tt.current), 0644); err != nil {
			t.Fatal(err)
		}
		if err := pet.Export(tt.in, dest); err != nil {
			t.Fatal(err)
		}
		d, err := os.ReadFile(dest)
		if err != nil {
			t.Fatal(err)
		}
		got := string(d)
		if got != tt.want {
			t.Errorf("got %v\nwant %v", got, tt.want)
		}
	}
}
