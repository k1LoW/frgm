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
		status, err := runFill(args)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		os.Exit(status)
	},
}

func runFill(args []string) (int, error) {
	exporter := frgm.New(config.Get("global.ignore").([]string))
	snippets, err := exporter.Load(srcPath)
	if err != nil {
		return 1, err
	}
	err = exporter.Export(snippets, srcPath)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func init() {
	config.Load()
	fillCmd.Flags().StringVarP(&srcPath, "src", "", config.Get("global.snippets_path").(string), "frgm snippets path")
	rootCmd.AddCommand(fillCmd)
}
