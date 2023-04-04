package frgm

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/k1LoW/frgm/format"
	"github.com/k1LoW/frgm/snippet"
	"github.com/karrick/godirwalk"
	gitignore "github.com/sabhiram/go-gitignore"
	"gopkg.in/jmhodges/yaml.v2"
)

var AllowExts = format.Exts{".yml", ".yaml"}

type Frgm struct {
	ignore []string
}

// New return new Frgm
func New(ignore []string) *Frgm {
	return &Frgm{
		ignore: ignore,
	}
}

func (f *Frgm) Load(src string) (snippet.Snippets, error) {
	snippets := snippet.Snippets{}
	if _, err := os.Lstat(src); err != nil {
		return snippets, err
	}
	i := gitignore.CompileIgnoreLines(f.ignore...)
	err := godirwalk.Walk(src, &godirwalk.Options{
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
			if !AllowExts.Contains(ext) {
				return nil
			}
			group := filepath.Base(fn[:len(fn)-len(ext)])
			set, err := f.LoadSet(file, group)
			if err != nil {
				return err
			}
			set.Snippets.AddLoadPath(path)
			snippets = append(snippets, set.Snippets...)
			return nil
		},
	})

	if err != nil {
		return snippets, err
	}

	return snippets, nil
}

func (f *Frgm) LoadSets(src string) ([]*snippet.SnippetSet, error) {
	sets := []*snippet.SnippetSet{}
	if _, err := os.Lstat(src); err != nil {
		return sets, err
	}
	i := gitignore.CompileIgnoreLines(f.ignore...)
	err := godirwalk.Walk(src, &godirwalk.Options{
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
			if !AllowExts.Contains(ext) {
				return nil
			}
			group := filepath.Base(fn[:len(fn)-len(ext)])
			set, err := f.LoadSet(file, group)
			if err != nil {
				return err
			}
			if len(set.Snippets) == 0 {
				return nil
			}
			set.LoadPath = path
			sets = append(sets, set)
			return nil
		},
	})

	if err != nil {
		return sets, err
	}

	return sets, nil
}

func (f *Frgm) LoadSet(in io.Reader, defaultGroup string) (*snippet.SnippetSet, error) {
	set := &snippet.SnippetSet{}
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return set, err
	}
	if !bytes.Contains(buf, []byte("snippets:")) {
		return set, nil
	}
	if err := yaml.Unmarshal(buf, set); err != nil {
		return set, err
	}
	for _, s := range set.Snippets {
		if s.Group == "" {
			if set.Group != "" {
				s.Group = set.Group
			} else {
				s.Group = defaultGroup
			}
		}
		if s.UID == "" {
			uid, err := genUID(s.Group, s.Name, s.Content)
			if err != nil {
				return set, err
			}
			s.UID = uid
		}
	}
	return set, nil
}

func (f *Frgm) Export(snippets snippet.Snippets, dest string) error {
	snippets.ClearLoadPath()
	sets, err := f.LoadSets(dest)
	if err != nil {
		return err
	}
L:
	for _, s := range snippets {
		// same UID
		for _, set := range sets {
			current, err := set.Snippets.FindByUID(s.UID)
			if err == nil {
				if s.Group != "" {
					current.Group = s.Group
				}
				current.Content = s.Content
				if s.Desc != "" {
					current.Desc = s.Desc
				}
				if len(s.Labels) > 0 {
					current.Labels = s.Labels
				}
				continue L
			}
		}
		if s.Group != "" {
			// same group ( SnippetSet )
			for _, set := range sets {
				if set.Group == s.Group {
					set.Snippets = append(set.Snippets, s)
					continue L
				}
			}
			// same group ( Snippet )
			for _, set := range sets {
				if set.Group != "" {
					continue
				}
				for _, c := range set.Snippets {
					if c.Group == s.Group {
						set.Snippets = append(set.Snippets, s)
						continue L
					}
				}
			}
		}
		// none
		sets[0].Snippets = append(sets[0].Snippets, s)
	}

	// set path
	paths := map[string]struct{}{}
	for _, set := range sets {
		if set.LoadPath != "" {
			if _, e := paths[set.LoadPath]; e {
				return fmt.Errorf("duplicate path: %s", set.LoadPath)
			}
			continue
		}
		if set.Group != "" {
			path := filepath.Join(dest, fmt.Sprintf("%s.yml", set.Group))
			if _, e := paths[path]; !e {
				paths[path] = struct{}{}
				set.LoadPath = path
				continue
			}
		}
		for {
			uid, err := genUID("", set.Snippets[0].Desc, set.Snippets[0].Content)
			if err != nil {
				return err
			}
			path := filepath.Join(dest, fmt.Sprintf("%s.yml", uid))
			if _, e := paths[path]; !e {
				set.LoadPath = path
				break
			}
		}
	}

	// export
	for _, set := range sets {
		path := set.LoadPath
		file, err := os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		encoder := yaml.NewEncoder(file)
		encoder.SetLineLength(-1)
		if err := encoder.Encode(set); err != nil {
			return err
		}
	}

	return nil
}

func (f *Frgm) Encode(out io.Writer, snippets snippet.Snippets) error {
	set := snippet.SnippetSet{
		Group:    "",
		Snippets: snippets,
	}
	encoder := yaml.NewEncoder(out)
	encoder.SetLineLength(-1)
	if err := encoder.Encode(set); err != nil {
		return err
	}
	return nil
}

func genUID(g, d, c string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, fmt.Sprintf("%s-%s-%s", g, d, c)); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	return fmt.Sprintf("frgm-%s", s[:12]), nil
}
