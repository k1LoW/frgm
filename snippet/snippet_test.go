package snippet

import "testing"

func TestValidate(t *testing.T) {
	snips := Snippets{
		New("uid", "group", "desc", "content", []string{}),
		New("uid", "group2", "desc2", "content2", []string{}),
	}
	if err := snips.Validate(); err == nil {
		t.Error("got nil want error")
	}
}
