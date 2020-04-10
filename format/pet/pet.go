package pet

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/k1LoW/frgm/format"
	"github.com/k1LoW/frgm/snippet"
)

var allowExts = format.Exts{".toml"}

type Snippet struct {
	Command     string   `toml:"command"`
	Description string   `toml:"description"`
	Output      string   `toml:"output"`
	Tag         []string `toml:"tag"`
}

type Snippets struct {
	Snippets []Snippet `toml:"snippets"`
}

type Pet struct{}

// New return new Pet
func New() *Pet {
	return &Pet{}
}

func (p *Pet) Load(src string) (snippet.Snippets, error) {
	snippets := snippet.Snippets{}
	if _, err := os.Lstat(src); err != nil {
		return snippets, err
	}
	file, err := os.Open(filepath.Clean(src))
	if err != nil {
		return snippets, err
	}
	fn := file.Name()
	ext := filepath.Ext(fn)
	if !allowExts.Contains(ext) {
		return snippets, fmt.Errorf("pet snippet file should be .toml: %s", src)
	}
	psnips, err := p.LoadFile(file)
	if err != nil {
		return snippets, err
	}
	for _, s := range psnips.Snippets {
		uid, err := genUID(s.Command)
		if err != nil {
			return snippets, err
		}
		snippets = append(snippets, snippet.New(uid, "default", s.Description, s.Command, s.Output, "", s.Tag))
	}
	return snippets, nil
}

func (p *Pet) LoadFile(in io.Reader) (Snippets, error) {
	var psnips Snippets
	if _, err := toml.DecodeReader(in, &psnips); err != nil {
		return psnips, err
	}
	return psnips, nil
}

func genUID(c string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, c); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	return fmt.Sprintf("pet-%s", s[:12]), nil
}

func (p *Pet) Export(snippets snippet.Snippets, dest string) error {
	psnips := Snippets{
		Snippets: []Snippet{},
	}
	current, err := p.Load(dest)
	if err != nil {
		return err
	}
	for _, s := range snippets {
		cs, err := current.FindByUID(s.UID)
		if err != nil {
			current = append(current, s)
			continue
		}
		cs.Name = s.Name
		cs.Desc = s.Desc
		cs.Output = s.Output
		cs.Labels = s.Labels
	}
	for _, s := range current {
		psnips.Snippets = append(psnips.Snippets, Snippet{
			Command:     s.Content,
			Description: s.Name,
			Output:      s.Output,
			Tag:         s.Labels,
		})
	}
	file, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(psnips); err != nil {
		return err
	}
	return nil
}
