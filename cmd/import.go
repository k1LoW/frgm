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

	"github.com/k1LoW/frgm/config"
	"github.com/k1LoW/frgm/format"
	"github.com/k1LoW/frgm/format/alfred"
	"github.com/k1LoW/frgm/format/frgm"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import other app snippets as frgm snippets",
	Long:  `Import other app snippets as frgm snippets.`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := runImport(args)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(status)
	},
}

func runImport(args []string) (int, error) {
	destPath = config.GetString("global.snippets_path")
	var (
		loader   format.Loader
		exporter format.Exporter
	)
	switch formatType {
	case "alfred":
		loader = alfred.New(config.GetStringSlice("global.ignore"))
	default:
		return 1, fmt.Errorf("unsupported format '%s'", formatType)
	}
	exporter = frgm.New(config.GetStringSlice("global.ignore"))

	snippets, err := loader.Load(srcPath)
	if err != nil {
		return 1, err
	}

	if err := snippets.Validate(); err != nil {
		return 1, err
	}

	err = exporter.Export(snippets, destPath)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&srcPath, "from", "f", "", "import snippets path")
	importCmd.Flags().StringVarP(&formatType, "format", "T", "alfred", "import snippets format of snippet")
}
