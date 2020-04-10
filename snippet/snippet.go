package snippet

import (
	"encoding/json"
	"fmt"
)

type Snippet struct {
	UID      string   `json:"uid,omitempty"`
	Group    string   `json:"group,omitempty"`
	Name     string   `json:"name"`
	Content  string   `json:"content"`
	Output   string   `json:"output,omitempty"`
	Desc     string   `json:"desc,omitempty"`
	Labels   []string `json:"labels,omitempty"`
	LoadPath string   `json:"-"`
}

func (s Snippet) String() string {
	j, _ := json.Marshal(s)
	return string(j)
}

// New return new Snippet
func New(u, g, n, c, o, d string, l []string) *Snippet {
	return &Snippet{
		UID:     u,
		Group:   g,
		Name:    n,
		Content: c,
		Output:  o,
		Desc:    d,
		Labels:  l,
	}
}

type Snippets []*Snippet

func (snips Snippets) AddLoadPath(path string) {
	for i := range snips {
		snips[i].LoadPath = path
	}
}

func (snips Snippets) ClearLoadPath() {
	for i := range snips {
		snips[i].LoadPath = ""
	}
}

func (snips Snippets) FindByUID(uid string) (*Snippet, error) {
	for _, s := range snips {
		if uid == s.UID {
			return s, nil
		}
	}
	return nil, fmt.Errorf("can not find snippet UID:%s", uid)
}

func (snips Snippets) Validate() error {
	uids := map[string]*Snippet{}

	// check for duplicate
	for _, s := range snips {
		if _, ok := uids[s.UID]; ok {
			return fmt.Errorf("duplicate UID\n%s\n---\n%s", uids[s.UID], s)
		}
		uids[s.UID] = s
	}

	return nil
}

type SnippetSet struct {
	Group    string   `yaml:"group,omitempty"`
	Snippets Snippets `yaml:"snippets"`
	LoadPath string   `json:"-"`
}
