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
	srcPath = config.Get("global.snippets_path").(string)
	loader := frgm.New(config.Get("global.ignore").([]string))
	snippets, err := loader.Load(srcPath)
	if err != nil {
		return 1, err
	}
	for _, s := range snippets {
		r := strings.NewReplacer(":uid", s.UID, ":group", s.Group, ":name", s.Name, ":content", s.Content, ":labels", fmt.Sprintf("%s", s.Labels), ":desc", s.Desc)
		fmt.Println(r.Replace(listFormat))
	}
	return 0, nil
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listFormat, "format", "", ":content # :name :labels", "list format")
}
