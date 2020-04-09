package format

import (
	"fmt"
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
		"label_join": func(l []string) string {
			return fmt.Sprintf("`%s`", strings.Join(l, "` `"))
		},
	}
}
