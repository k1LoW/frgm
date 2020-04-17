package history

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
			in:   ``,
			want: snippet.Snippets{},
		},
		{
			in: `
`,
			want: snippet.Snippets{},
		},
		{
			in: ` 8  man history
 9  git checkout -b format-history
10  emacs format/
11* go test ./format/history
12  ls # exec ls
`,
			want: snippet.Snippets{
				&snippet.Snippet{
					UID:     "history-bd725c52840f",
					Group:   "default",
					Name:    "man history",
					Content: "man history",
					Labels:  []string{},
				},
				&snippet.Snippet{
					UID:     "history-709a48a0dc97",
					Group:   "default",
					Name:    "git checkout -b format-history",
					Content: "git checkout -b format-history",
					Labels:  []string{},
				},
				&snippet.Snippet{
					UID:     "history-fd42cee8fd75",
					Group:   "default",
					Name:    "emacs format/",
					Content: "emacs format/",
					Labels:  []string{},
				},
				&snippet.Snippet{
					UID:     "history-00b0136f1cf7",
					Group:   "default",
					Name:    "go test ./format/history",
					Content: "go test ./format/history",
					Labels:  []string{},
				},
				&snippet.Snippet{
					UID:     "history-a3a97692f42a",
					Group:   "default",
					Name:    "exec ls",
					Content: "ls",
					Labels:  []string{},
				},
			},
		},
	}

	for _, tt := range tests {
		in := bytes.NewBufferString(tt.in)
		history := New()
		got, err := history.Decode(in, "default")
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Errorf("got\n%v\nwant\n%v", got, tt.want)
		}
	}
}
