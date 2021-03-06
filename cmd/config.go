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

	"github.com/k1LoW/frgm/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get and set frgm config",
	Long:  `Get and set frgm config.`,
	RunE:  runConfig,
}

func runConfig(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		for k, v := range viper.AllSettings() {
			switch v := v.(type) {
			case map[string]interface{}:
				for kk, vv := range v {
					cmd.Printf("%s.%s=%v\n", k, kk, vv)
				}
			default:
				cmd.Printf("%s=%v\n", k, v)
			}
		}
	case len(args) == 1:
		cmd.Printf("%v\n", viper.Get(args[0]))
	case len(args) == 2:
		if err := config.Set(args[0], args[1]); err != nil {
			return err
		}
		if err := config.Save(); err != nil {
			return err
		}
	default:
		return errors.New("invalid arguments")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)
}
