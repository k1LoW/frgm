package md

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/k1LoW/frgm/format"
	"github.com/k1LoW/frgm/snippet"
)

var allowExts = format.Exts{".md"}

//go:embed templates/*
var tmpl embed.FS

type Md struct {
	ignore []string
	tmpl   embed.FS
}

// New return *Md
func New(ignore []string) *Md {
	return &Md{
		ignore: ignore,
		tmpl:   tmpl,
	}
}

func (m *Md) Export(snippets snippet.Snippets, dest string) error {
	f, err := os.Stat(dest)
	if err == nil && f.IsDir() {
		return fmt.Errorf("%s is directory", dest)
	}
	ext := filepath.Ext(dest)
	if !allowExts.Contains(ext) {
		return fmt.Errorf("%s is not .md file", dest)
	}
	snippets.ClearLoadPath()
	sets := []*snippet.SnippetSet{}
	groups := map[string]*snippet.SnippetSet{}

	for _, s := range snippets {
		set, ok := groups[s.Group]
		if !ok {
			set = &snippet.SnippetSet{
				Group:    s.Group,
				Snippets: snippet.Snippets{},
			}
			groups[s.Group] = set
			sets = append(sets, set)
		}
		set.Snippets = append(set.Snippets, s)
	}

	path := dest
	file, err := os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if err := m.Write(file, sets); err != nil {
		return err
	}

	return nil
}

func (m *Md) Write(wr io.Writer, sets []*snippet.SnippetSet) error {
	tb, err := m.tmpl.ReadFile("templates/snippets.md.tmpl")
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("snippets").Funcs(format.Funcs()).Parse(string(tb)))
	tmplData := map[string]interface{}{
		"Sets": sets,
	}
	if err != nil {
		return err
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}
