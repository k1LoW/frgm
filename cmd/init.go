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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Songmu/prompter"
	"github.com/k1LoW/frgm/config"
	"github.com/k1LoW/frgm/format/frgm"
	"github.com/spf13/cobra"
)

var force bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize frgm",
	Long:  `Initialize frgm.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runInit(cmd, args)
		if err != nil {
			printErrln(cmd, err)
			os.Exit(1)
		}
	},
}

func runInit(cmd *cobra.Command, args []string) error {
	path := config.GetString("global.snippets_path")
	if !force {
		path = prompter.Prompt("Enter snippet path (dir or .yml file) [global.snippets_path]", path)
	}
	if err := config.Set("global.snippets_path", path); err != nil {
		return err
	}
	_, err := os.Stat(path)
	if err != nil {
		ext := filepath.Ext(path)
		if frgm.AllowExts.Contains(ext) {
			var yn bool
			if force {
				yn = true
			} else {
				yn = prompter.YN(fmt.Sprintf("Create new snippet file (%s) ?", path), true)
			}
			if yn {
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path, []byte("---\nsnippets: []\n"), 0600); err != nil {
					return err
				}
			}
		} else {
			var yn bool
			if force {
				yn = true
			} else {
				yn = prompter.YN(fmt.Sprintf("Create snippets directory (%s)?", path), true)
			}
			if yn {
				if err := os.MkdirAll(path, 0700); err != nil {
					return err
				}
			}
		}
		cmd.Printf("Create %s\n", path)
	}

	if err := config.Save(); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "force init")
}
