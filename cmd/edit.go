/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
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
		status, err := runEdit(args)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(status)
	},
}

func runEdit(args []string) (int, error) {
	e := os.Getenv("EDITOR")
	if e == "" {
		return 1, errors.New("$EDITOR is not set")
	}
	tty, err := tty.Open()
	if err != nil {
		return 1, err
	}
	defer tty.Close()
	c := exec.Command(e, config.GetString("global.snippets_path")) // #nosec
	c.Stdin = tty.Input()
	c.Stdout = tty.Output()
	c.Stderr = tty.Output()
	if err := c.Run(); err != nil {
		return 1, err
	}
	return 0, nil
}

func init() {
	rootCmd.AddCommand(editCmd)
}
