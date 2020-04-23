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
	"github.com/k1LoW/frgm/format/md"
	"github.com/k1LoW/frgm/format/pet"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export frgm snippets as other app snippets",
	Long:  `Export frgm snippets as othre app snippets.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runExport(args)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

func runExport(args []string) error {
	var (
		loader   format.Loader
		exporter format.Exporter
	)
	loader = frgm.New(config.GetStringSlice("global.ignore"))

	switch formatType {
	case "alfred":
		exporter = alfred.New(config.GetStringSlice("global.ignore"))
	case "pet":
		exporter = pet.New()
	case "md":
		exporter = md.New(config.GetStringSlice("global.ignore"))
	default:
		return fmt.Errorf("unsupported format '%s'", formatType)
	}

	if srcPath == "" {
		srcPath = config.GetString("global.snippets_path")
	}

	snippets, err := loader.Load(srcPath)
	if err != nil {
		return err
	}

	if err := snippets.Validate(); err != nil {
		return err
	}

	err = exporter.Export(snippets, destPath)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	config.Load()
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&srcPath, "from", "f", config.GetString("global.snippets_path"), "frgm snippets path")
	exportCmd.Flags().StringVarP(&destPath, "to", "t", "", "export snippets path")
	exportCmd.Flags().StringVarP(&formatType, "format", "T", "alfred", "export format of snippet")
}
