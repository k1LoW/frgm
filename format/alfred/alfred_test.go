package alfred

import (
	"bytes"
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
			in: `{
  "alfredsnippet": {
    "uid": "frgm-46c29e119523",
    "name": "ping",
    "snippet": "ping 8.8.8.8",
    "keyword": "network google"
  }
}
`,
			want: snippet.Snippets{
				&snippet.Snippet{
					UID:     "frgm-46c29e119523",
					Group:   "default",
					Content: "ping 8.8.8.8",
					Name:    "ping",
					Output:  "",
					Labels:  []string{"network", "google"},
				},
			},
		},
	}
	for _, tt := range tests {
		in := bytes.NewBufferString(tt.in)
		a := New([]string{})
		got, err := a.Decode(in, "default")
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Errorf("got %v\nwant %v", got, tt.want)
		}
	}
}
