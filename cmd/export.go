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
	"github.com/k1LoW/frgm/drivers"
	"github.com/k1LoW/frgm/drivers/alfred"
	"github.com/k1LoW/frgm/drivers/frgm"
	"github.com/spf13/cobra"
)

var (
	srcPath  string
	destPath string
	format   string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export frgm snippets as other app snippets",
	Long:  `Export frgm snippets as othre app snippets.`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := runExport(args)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(status)
	},
}

func runExport(args []string) (int, error) {
	var (
		loader   drivers.Loader
		exporter drivers.Exporter
	)
	loader = frgm.New(config.Get("global.ignore").([]string))

	switch format {
	case "alfred":
		exporter = alfred.New(config.Get("global.ignore").([]string))
	default:
		return 1, fmt.Errorf("unsupported driver '%s'", format)
	}

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
	config.Load()
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&srcPath, "src", "", config.Get("global.snippets_path").(string), "frgm snippets dir")
	exportCmd.Flags().StringVarP(&destPath, "dest", "", "", "export destination")
	exportCmd.Flags().StringVarP(&format, "format", "T", "alfred", "export format of snippet")
}
