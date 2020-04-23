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
	"os"

	"github.com/k1LoW/frgm/config"
	"github.com/k1LoW/frgm/format/frgm"
	"github.com/spf13/cobra"
)

// fillCmd represents the fill command
var fillCmd = &cobra.Command{
	Use:   "fill",
	Short: "Fill in the blanks in current snippets",
	Long:  `Fill in the blanks in current snippets.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runFill(args)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

func runFill(args []string) error {
	srcPath = config.GetString("global.snippets_path")
	exporter := frgm.New(config.GetStringSlice("global.ignore"))
	snippets, err := exporter.Load(srcPath)
	if err != nil {
		return err
	}
	err = exporter.Export(snippets, srcPath)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(fillCmd)
}
