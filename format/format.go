package format

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/k1LoW/frgm/snippet"
)

type Loader interface {
	Load(src string) (snippet.Snippets, error)
}

type Exporter interface {
	Export(snippets snippet.Snippets, dest string) error
}

type Decoder interface {
	Decode(in io.Reader, group string) (snippet.Snippets, error)
}

type Encoder interface {
	Encode(out io.Writer, snippets snippet.Snippets) error
}

type Exts []string

func (exts Exts) Contains(t string) bool {
	for _, e := range exts {
		if e == t {
			return true
		}
	}
	return false
}

func Funcs() map[string]interface{} {
	return template.FuncMap{
		"nl2mdnl": func(text string) string {
			r := strings.NewReplacer("\r\n", "  \n", "\n", "  \n", "\r", "  \n")
			return r.Replace(text)
		},
		"label_join": func(l []string) string {
			return fmt.Sprintf("`%s`", strings.Join(l, "` `"))
		},
	}
}
