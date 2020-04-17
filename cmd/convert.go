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
	"errors"
	"fmt"
	"os"

	"github.com/k1LoW/frgm/config"
	"github.com/k1LoW/frgm/format"
	"github.com/k1LoW/frgm/format/alfred"
	"github.com/k1LoW/frgm/format/frgm"
	"github.com/k1LoW/frgm/format/history"
	"github.com/k1LoW/frgm/format/pet"
	"github.com/spf13/cobra"
)

var group string

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert other app snippet format to frgm snippet format using STDIN/STDOUT",
	Long:  `Convert other app snippet format to frgm snippet format using STDIN/STDOUT.`,
	Args: func(cmd *cobra.Command, args []string) error {
		fi, err := os.Stdin.Stat()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			return errors.New("ghput need STDIN. Please use pipe")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		status, err := runConvert(args)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(status)
	},
}

func runConvert(args []string) (int, error) {
	var (
		decoder format.Decoder
		encoder format.Encoder
	)
	encoder = frgm.New(config.GetStringSlice("global.ignore"))

	switch formatType {
	case "alfred":
		decoder = alfred.New(config.GetStringSlice("global.ignore"))
	case "pet":
		decoder = pet.New()
	case "history":
		decoder = history.New()
	default:
		return 1, fmt.Errorf("unsupported format '%s'", formatType)
	}

	snippets, err := decoder.Decode(os.Stdin, group)
	if err != nil {
		return 1, err
	}
	if err := encoder.Encode(os.Stdout, snippets); err != nil {
		return 1, err
	}

	return 0, nil
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&group, "group", "g", "default", "group of snippets")
	convertCmd.Flags().StringVarP(&formatType, "format", "T", "alfred", "format of STDIN snippet")
}
