package frgm

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/k1LoW/frgm/snippet"
)

type Frgm struct{}

// New return new Frgm
func New() *Frgm {
	return &Frgm{}
}

func (f *Frgm) Load(src string) ([]snippet.Snippet, error) {
	snippets := []snippet.Snippet{}
	if _, err := os.Lstat(src); err != nil {
		return snippets, err
	}

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(filepath.Clean(path))
		if err != nil {
			return err
		}
		fn := file.Name()
		group := filepath.Base(fn[:len(fn)-len(filepath.Ext(fn))])
		s, err := f.LoadSet(file, group)
		if err != nil {
			return err
		}
		snippets = append(snippets, s...)
		return nil
	})

	if err != nil {
		return snippets, err
	}

	return snippets, nil
}

func (f *Frgm) LoadSet(in io.Reader, defaultGroup string) ([]snippet.Snippet, error) {
	snippets := []snippet.Snippet{}
	set := &snippet.SnippetSet{}
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return snippets, err
	}
	err = yaml.Unmarshal(buf, set)
	if err != nil {
		return snippets, err
	}
	for _, s := range set.Snippets {
		if s.Group == "" {
			s.Group = defaultGroup
		}
		if s.UID == "" {
			uid, err := genUID(s.Group, s.Desc, s.Content)
			if err != nil {
				return snippets, err
			}
			s.UID = uid
		}
		snippets = append(snippets, s)
	}
	return snippets, nil
}

func genUID(g, d, c string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, fmt.Sprintf("%s-%s-%s", g, d, c)); err != nil {
		return "", err
	}
	return fmt.Sprintf("frgm-%x", h.Sum(nil)), nil
}
