package alfred

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/k1LoW/frgm/format"
	"github.com/k1LoW/frgm/snippet"
	"github.com/karrick/godirwalk"
	gitignore "github.com/sabhiram/go-gitignore"
)

const DefaultDest = "~/Library/Application Support/Alfred/Alfred.alfredpreferences/snippets"

var allowExts = format.Exts{".json"}

type Alfred struct {
	ignore []string
}

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
func New(ignore []string) *Alfred {
	return &Alfred{
		ignore: ignore,
	}
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

func (a *Alfred) Load(src string) (snippet.Snippets, error) {
	snippets := snippet.Snippets{}
	d, err := os.Stat(src)
	if err != nil {
		return snippets, err
	}
	if !d.IsDir() {
		return snippets, fmt.Errorf("%s is not directory", src)
	}
	i, err := gitignore.CompileIgnoreLines(a.ignore...)
	if err != nil {
		return snippets, err
	}
	err = godirwalk.Walk(src, &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			b, err := de.IsDirOrSymlinkToDir()
			if err != nil {
				return err
			}
			if b {
				return nil
			}
			if i.MatchesPath(path) {
				return nil
			}
			file, err := os.Open(filepath.Clean(path))
			if err != nil {
				return err
			}
			fn := file.Name()
			ext := filepath.Ext(fn)
			if !allowExts.Contains(ext) {
				return nil
			}
			group := filepath.Base(filepath.Dir(path))
			s, err := a.LoadOne(file, group)
			if err != nil {
				return err
			}
			s.LoadPath = path
			snippets = append(snippets, s)
			return nil
		},
	})

	if err != nil {
		return snippets, err
	}

	return snippets, nil
}

func (a *Alfred) LoadOne(in io.Reader, group string) (*snippet.Snippet, error) {
	as := &Snippet{}
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, as); err != nil {
		return nil, err
	}
	return snippet.New(as.AlfredSnippet.UID, group, as.AlfredSnippet.Name, as.AlfredSnippet.Snippet, "", strings.Split(as.AlfredSnippet.Keyword, " ")), nil
}

func (a *Alfred) Export(snippets snippet.Snippets, dest string) error {
	snippets.ClearLoadPath()
	current, err := a.Load(dest)
	if err != nil {
		return err
	}
	for _, s := range snippets {
		path := filepath.Join(dest, s.Group, fmt.Sprintf("%s.json", s.UID))
		cs, err := current.FindByUID(s.UID)
		if err == nil {
			if cs.Group != s.Group {
				return fmt.Errorf("group is different for existing snippet and new snippet. Delete or move the existing snippet or match the group of the new snippet\n%s\n---\n%s", cs, s)
			}
			path = cs.LoadPath
		}
		as := NewSnippet(s.UID, s.Name, s.Content, strings.Join(s.Labels, " "))
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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
