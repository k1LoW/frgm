package drivers

import "github.com/k1LoW/frgm/snippet"

type Loader interface {
	Load(src string) (snippet.Snippets, error)
}

type Exporter interface {
	Export(snippets snippet.Snippets, dest string) error
}
