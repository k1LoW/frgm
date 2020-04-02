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
	destPath = config.Get("global.snippets_path").(string)
	var (
		loader   format.Loader
		exporter format.Exporter
	)
	switch formatType {
	case "alfred":
		loader = alfred.New(config.Get("global.ignore").([]string))
	default:
		return 1, fmt.Errorf("unsupported format '%s'", formatType)
	}
	exporter = frgm.New(config.Get("global.ignore").([]string))

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
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&srcPath, "from", "f", "", "import snippets path")
	importCmd.Flags().StringVarP(&formatType, "format", "T", "alfred", "import snippets format of snippet")
}
