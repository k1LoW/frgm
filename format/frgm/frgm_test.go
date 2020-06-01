package frgm

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k1LoW/frgm/snippet"
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

	frgm := New([]string{})
	for _, tt := range tests {
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

func TestEncode(t *testing.T) {
	tests := []struct {
		snippets snippet.Snippets
		want     string
	}{
		{
			snippets: snippet.Snippets{},
			want: `snippets: []
`,
		},
		{
			snippets: snippet.Snippets{
				snippet.New("frgm-de1fbe2f7573", "my-group", "hello world", "echo hello world", "hello world", "show following\nhello world\n", []string{}),
			},
			want: `snippets:
- uid: frgm-de1fbe2f7573
  group: my-group
  name: hello world
  desc: |
    show following
    hello world
  content: echo hello world
  output: hello world
`,
		},
		{
			snippets: snippet.Snippets{
				snippet.New("frgm-de1fbe2f7573", "my-group", "count list", `cat list | sort -nk1`, "", "", []string{"a-b", "c"}),
			},
			want: `snippets:
- uid: frgm-de1fbe2f7573
  group: my-group
  name: count list
  content: cat list | sort -nk1
  labels:
  - a-b
  - c
`,
		},
		{
			snippets: snippet.Snippets{
				snippet.New("frgm-de1fbe2f7573", "my-group", "count access.log", `cat /var/log/nginx/access.log | awk '{ if (($1 == "domain:example.com") && ($8 ~ /^status:[2345]/)) { arr[$8]++; x++ } } END { for (i in arr) { printf "%s\t%.2f%\t%s / %s\n", i, arr[i]/x*100, arr[i], x } }' | sort -nk1`, "", "", []string{"a-b", "c"}),
			},
			want: `snippets:
- uid: frgm-de1fbe2f7573
  group: my-group
  name: count access.log
  content: cat /var/log/nginx/access.log | awk '{ if (($1 == "domain:example.com") && ($8 ~ /^status:[2345]/)) { arr[$8]++; x++ } } END { for (i in arr) { printf "%s\t%.2f%\t%s / %s\n", i, arr[i]/x*100, arr[i], x } }' | sort -nk1
  labels:
  - a-b
  - c
`,
		},
	}

	frgm := New([]string{})
	for _, tt := range tests {
		out := new(bytes.Buffer)
		if err := frgm.Encode(out, tt.snippets); err != nil {
			t.Fatal(err)
		}
		got := out.String()
		if got != tt.want {
			t.Errorf("got %v\nwant %v", got, tt.want)
		}
		set, err := frgm.LoadSet(out, "default")
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(set.Snippets, tt.snippets, nil); diff != "" {
			t.Errorf("%s", diff)
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
