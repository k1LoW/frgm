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
		frgm := New([]string{})
		got, err := frgm.LoadSet(in, "default-group")
		if err != nil {
			t.Fatal(err)
		}
		if want := 1; len(got.Snippets) != want {
			t.Errorf("got %v\nwant %v", len(got.Snippets), want)
		}
		if want := tt.wantGroup; got.Snippets[0].Group != want {
			t.Errorf("got %v\nwant %v", got.Snippets[0].Group, want)
		}
		if ok, _ := regexp.MatchString(tt.wantUIDMatch, got.Snippets[0].UID); !ok {
			t.Errorf("got %v\nwant match %v", got.Snippets[0].UID, tt.wantUIDMatch)
		}
	}
}

func TestGenUID(t *testing.T) {
	got, err := genUID("", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if want := "frgm-d8156bae0c42"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
