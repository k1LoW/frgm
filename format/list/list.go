package list

import (
	"fmt"
	"io"
	"strings"

	"github.com/k1LoW/frgm/snippet"
)

const DefaultFormat = ":content # :name [:group :labels]"

var nlr = strings.NewReplacer("\\\n", "")

type List struct {
	format string
}

// New return new List
func New(format string) *List {
	return &List{
		format: format,
	}
}

func (f *List) Encode(out io.Writer, snippets snippet.Snippets) error {
	for _, s := range snippets {
		r := strings.NewReplacer(":uid", s.UID, ":group", s.Group, ":name", s.Name, ":content", nlr.Replace(s.Content), ":labels", strings.Join(s.Labels, " "), ":desc", s.Desc, ":output", s.Output)
		_, err := fmt.Fprintln(out, r.Replace(f.format))
		if err != nil {
			return err
		}
	}
	return nil
}
