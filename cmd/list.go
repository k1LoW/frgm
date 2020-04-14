/*
Copyright © 2020 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/k1LoW/frgm/config"
	"github.com/k1LoW/frgm/format/frgm"
	"github.com/spf13/cobra"
)

var listFormat string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List snippets",
	Long:  `List snippets.`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := runList(args)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(status)
	},
}

func runList(args []string) (int, error) {
	srcPath = config.GetString("global.snippets_path")
	loader := frgm.New(config.GetStringSlice("global.ignore"))
	snippets, err := loader.Load(srcPath)
	if err != nil {
		return 1, err
	}
	for _, s := range snippets {
		r := strings.NewReplacer(":uid", s.UID, ":group", s.Group, ":name", s.Name, ":content", s.Content, ":labels", strings.Join(s.Labels, " "), ":desc", s.Desc)
		fmt.Println(r.Replace(listFormat))
	}
	return 0, nil
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listFormat, "format", "", ":content # :name [:group :labels]", "list format")
}
