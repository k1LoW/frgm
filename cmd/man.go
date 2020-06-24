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
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/k1LoW/frgm/config"
	"github.com/k1LoW/frgm/format/frgm"
	"github.com/labstack/gommon/color"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// manCmd represents the man command
var manCmd = &cobra.Command{
	Use:   "man [UID]",
	Short: "Show command man",
	Long:  `Show command man.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return nil
		}
		if isatty.IsTerminal(os.Stdin.Fd()) {
			return errors.New("accepts 1 arg(s), received 0")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runMan(args)
		if err != nil {
			printFatalln(cmd, err)
		}
	},
}

func runMan(args []string) error {
	var uid string
	if len(args) == 1 {
		uid = args[0]
	} else {
		stdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		uid = strings.TrimSpace(string(stdin))
	}
	loader := frgm.New(config.GetStringSlice("global.ignore"))

	if srcPath == "" {
		srcPath = config.GetString("global.snippets_path")
	}
	snippets, err := loader.Load(srcPath)
	if err != nil {
		return err
	}
	s, err := snippets.FindByUID(uid)
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New("man").Funcs(funcs()).Parse(`{{ "NAME" | bold }}
       {{ .snippet.Name }}

{{ "CONTENT" | bold }}
       {{ .snippet.Content | nlindent | bold }}

{{ if ne .snippet.Output "" }}{{ "OUTPUT" | bold }}
       {{ .snippet.Output | nlindent }}

{{ end }}{{ if ne .snippet.Desc "" }}{{ "DESCRIPTION" | bold }}
       {{ .snippet.Desc | nlindent }}

{{ end }}{{ if ne .snippet.Group "" }}{{ "GROUP" | bold }}
       {{ .snippet.Group }}

{{ end }}{{ if gt (len .snippet.Labels) 0 }}{{ "LABELS" | bold }}
       {{ .snippet.Labels | labels}}

{{ end }}{{ if ne .snippet.LoadPath "" }}{{ "LOAD PATH" | bold }}
       {{ .snippet.LoadPath }}
{{ end }}
`))

	templateData := map[string]interface{}{
		"snippet": s,
	}
	err = tmpl.Execute(os.Stdout, templateData)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func funcs() map[string]interface{} {
	return map[string]interface{}{
		"nlindent": func(text string) string {
			r := strings.NewReplacer("\n", "\n       ", "\r\n", "\r\n       ", "\r", "\r       ")
			return r.Replace(text)
		},
		"bold": func(text string) string {
			return color.Bold(text)
		},
		"labels": func(labels []string) string {
			return strings.Join(labels, ", ")
		},
	}
}

func init() {
	config.Load()
	rootCmd.AddCommand(manCmd)
	manCmd.Flags().StringVarP(&srcPath, "from", "f", config.GetString("global.snippets_path"), "frgm snippets path")
}
