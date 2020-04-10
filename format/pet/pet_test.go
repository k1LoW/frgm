package pet

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
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
