package frgm

import (
	"bytes"
	"regexp"
	"testing"
)

func TestLoadSet(t *testing.T) {
	tests := []struct {
		in           string
		wantGroup    string
		wantUIDMatch string
	}{
		{
			in: `---
snippets:
  -
    name: hello world
    content: echo hello world
    label:
      - test
      - echo
`,
			wantGroup:    "default-group",
			wantUIDMatch: "^frgm",
		},
		{
			in: `---
snippets:
  -
    name: hello world
    content: echo hello world
    group: my-group
    label:
      - test
      - echo
`,
			wantGroup:    "my-group",
			wantUIDMatch: "^frgm",
		},
		{
			in: `---
snippets:
  -
    name: hello world
    content: echo hello world
    uid: 1234-5678-90
    label:
      - test
      - echo
`,
			wantGroup:    "default-group",
			wantUIDMatch: "^1234-5678-90$",
		},
	}

	for _, tt := range tests {
		in := bytes.NewBufferString(tt.in)
		frgm := New()
		got, err := frgm.LoadSet(in, "default-group")
		if err != nil {
			t.Fatal(err)
		}
		if want := 1; len(got) != want {
			t.Errorf("got %v\nwant %v", len(got), want)
		}
		if want := tt.wantGroup; got[0].Group != want {
			t.Errorf("got %v\nwant %v", got[0].Group, want)
		}
		if ok, _ := regexp.MatchString(tt.wantUIDMatch, got[0].UID); !ok {
			t.Errorf("got %v\nwant match %v", got[0].UID, tt.wantUIDMatch)
		}
	}
}
