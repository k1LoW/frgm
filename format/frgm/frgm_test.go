package frgm

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
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
    labels:
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
    labels:
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
    labels:
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

func TestFill(t *testing.T) {
	tests := []struct {
		in    string
		group string
		want  string
	}{
		{
			in: `---
snippets:
  -
    name: hello world
    content: echo hello world
`,
			group: "my-group",
			want: `---
snippets:
  -
    uid: frgm-de1fbe2f7573
    group: my-group
    name: hello world
    content: echo hello world
`,
		},
		{
			in: `---
snippets:
  -
    name: hello world
    content: echo hello world
    group: other-group
`,
			group: "my-group",
			want: `---
snippets:
  -
    uid: frgm-526b6b7787e4
    group: other-group
    name: hello world
    content: echo hello world
`,
		},
	}

	for _, tt := range tests {
		frgm := New([]string{})
		in := bytes.NewBufferString(tt.in)
		got, err := frgm.LoadSet(in, tt.group)
		if err != nil {
			t.Fatal(err)
		}

		out := bytes.NewBufferString(tt.want)
		want, err := frgm.LoadSet(out, tt.group)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(got, want, nil); diff != "" {
			t.Errorf("got %v\nwant %v", got, want)
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
