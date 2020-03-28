package snippet

type Snippet struct {
	UID     string   `yaml:"uid,omitempty"`
	Group   string   `yaml:"group,omitempty"`
	Desc    string   `yaml:"desc"`
	Content string   `yaml:"content"`
	Labels  []string `yaml:"labels,omitempty"`
}

type SnippetSet struct {
	Snippets []Snippet `yaml:"snippets"`
}
