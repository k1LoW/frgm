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
	"os"
	"os/exec"

	"github.com/k1LoW/frgm/config"
	"github.com/mattn/go-tty"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit global.snippets_path using $EDITOR",
	Long:  `Edit global.snippets_path using $EDITOR.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runEdit(args)
		if err != nil {
			printErrln(cmd, err)
			os.Exit(1)
		}
	},
}

func runEdit(args []string) error {
	e := os.Getenv("EDITOR")
	if e == "" {
		return errors.New("$EDITOR is not set")
	}
	tty, err := tty.Open()
	if err != nil {
		return err
	}
	defer tty.Close()
	c := exec.Command(e, config.GetString("global.snippets_path")) // #nosec
	c.Stdin = tty.Input()
	c.Stdout = tty.Output()
	c.Stderr = tty.Output()
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(editCmd)
}
