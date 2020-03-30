package snippet

import "testing"

func TestValidate(t *testing.T) {
	snips := Snippets{
		New("uid", "group", "name", "content", "desc", []string{}),
		New("uid", "group2", "name2", "content2", "desc2", []string{}),
	}
	if err := snips.Validate(); err == nil {
		t.Error("got nil want error")
	}
}
