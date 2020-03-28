package drivers

import "github.com/k1LoW/frgm/snippet"

type Loader interface {
	Load(src string) ([]snippet.Snippet, error)
}

type Exporter interface {
	Export(snippets []snippet.Snippet, dest string) error
}
