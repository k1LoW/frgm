package list

import (
	"bytes"
	"testing"

	"github.com/k1LoW/frgm/snippet"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		format   string
		snippets snippet.Snippets
		want     string
	}{
		{
			format:   DefaultFormat,
			snippets: snippet.Snippets{},
			want:     ``,
		},
		{
			format: DefaultFormat,
			snippets: snippet.Snippets{
				snippet.New("frgm-de1fbe2f7573", "my-group", "hello world", "echo hello world", "hello world", "show following\nhello world\n", []string{"a", "b"}),
			},
			want: `echo hello world # hello world [my-group a b]
`,
		},
		{
			format: ":content",
			snippets: snippet.Snippets{
				snippet.New("frgm-de1fbe2f7573", "my-group", "hello world", "echo hello world", "hello world", "show following\nhello world\n", []string{"a", "b"}),
			},
			want: `echo hello world
`,
		},
		{
			format: DefaultFormat,
			snippets: snippet.Snippets{
				snippet.New("frgm-de1fbe2f7573", "my-group", "hello world", `echo hello \
world`, "hello world", "show following\nhello world\n", []string{"a", "b"}),
			},
			want: `echo hello world # hello world [my-group a b]
`,
		},
	}

	for _, tt := range tests {
		list := New(tt.format)
		out := new(bytes.Buffer)
		if err := list.Encode(out, tt.snippets); err != nil {
			t.Fatal(err)
		}
		got := out.String()
		if got != tt.want {
			t.Errorf("got %v\nwant %v", got, tt.want)
		}
	}
}
