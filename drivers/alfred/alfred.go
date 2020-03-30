package alfred

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/k1LoW/frgm/snippet"
)

const DefaultDest = "~/Library/Application Support/Alfred/Alfred.alfredpreferences/snippets"

type Alfred struct{}

type Snippet struct {
	AlfredSnippet AlfredSnippet `json:"alfredsnippet"`
}

type AlfredSnippet struct {
	UID     string `json:"uid"`
	Name    string `json:"name"`
	Snippet string `json:"snippet"`
	Keyword string `json:"keyword"`
}

// New return *Alfred
func New() *Alfred {
	return &Alfred{}
}

func NewSnippet(u, n, s, k string) Snippet {
	return Snippet{
		AlfredSnippet: AlfredSnippet{
			UID:     u,
			Name:    n,
			Snippet: s,
			Keyword: k,
		},
	}
}

func (a *Alfred) Export(snippets snippet.Snippets, dest string) error {
	d, err := os.Stat(dest)
	if err != nil {
		return err
	}
	if !d.IsDir() {
		return fmt.Errorf("%s is not directory", dest)
	}
	for _, s := range snippets {
		as := NewSnippet(s.UID, s.Name, s.Content, strings.Join(s.Labels, " "))
		if err := os.MkdirAll(filepath.Join(dest, s.Group), 0700); err != nil {
			return err
		}
		file, err := os.OpenFile(filepath.Join(dest, s.Group, fmt.Sprintf("%s.json", s.UID)), os.O_WRONLY|os.O_CREATE, 0600) // #nosec
		if err != nil {
			return err
		}
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(as); err != nil {
			return err
		}
	}
	return nil
}
