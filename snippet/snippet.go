package snippet

import (
	"encoding/json"
	"fmt"
)

type Snippet struct {
	UID     string   `json:"uid,omitempty"`
	Group   string   `json:"group,omitempty"`
	Name    string   `json:"name"`
	Content string   `json:"content"`
	Labels  []string `json:"labels,omitempty"`
}

func (s Snippet) String() string {
	j, _ := json.Marshal(s)
	return string(j)
}

// New return new Snippet
func New(u, g, n, c string, l []string) Snippet {
	return Snippet{
		UID:     u,
		Group:   g,
		Name:    n,
		Content: c,
		Labels:  l,
	}
}

type Snippets []Snippet

func (snips Snippets) Validate() error {
	uids := map[string]Snippet{}

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
	Snippets Snippets `yaml:"snippets"`
}
