package md

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/k1LoW/frgm/snippet"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		sets       []*snippet.SnippetSet
		goldenPath string
	}{
		{
			[]*snippet.SnippetSet{
				&snippet.SnippetSet{
					Group: "group-a",
					Snippets: snippet.Snippets{
						snippet.New("uid-a-1", "group-a", "a-1", "echo a-1", "desc a-1", []string{}),
						snippet.New("uid-a-2", "group-a", "a-2", "echo a-2", "", []string{}),
					},
				},
				&snippet.SnippetSet{
					Group: "group-b",
					Snippets: snippet.Snippets{
						snippet.New("uid-b-1", "group-b", "b-1", "echo b-1", "desc b-1\ndesc b-1", []string{}),
					},
				},
			},
			"md_write_1.md.golden",
		},
	}

	for _, tt := range tests {
		md := New([]string{})
		got := new(bytes.Buffer)
		if err := md.Write(got, tt.sets); err != nil {
			t.Fatal(err)
		}
		want, err := ioutil.ReadFile(filepath.Join(testdataDir(), tt.goldenPath))
		if err != nil {
			t.Fatal(err)
		}
		if got.String() != string(want) {
			t.Errorf("got\n%s\nwant\n%s", got.Bytes(), want)
		}
	}
}

func testdataDir() string {
	wd, _ := os.Getwd()
	dir, _ := filepath.Abs(filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata"))
	return dir
}
