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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/k1LoW/frgm/config"
	"github.com/spf13/cobra"
	giturl "github.com/whilp/git-urls"
	"github.com/x-motemen/ghq/cmdutil"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage snippets repositories",
	Long:  `Manage snippets repositories.`,
}

var repoAddCmd = &cobra.Command{
	Use:   "add [REPO_URL]",
	Short: "Add snippets repository to 'global.snippets_path'",
	Long:  `Add snippets repository 'global.snippets_path'.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]
		rootPath := config.GetString("global.snippets_path")
		_, err := exec.LookPath("ghq")
		if err == nil && prompter.YN("Do you use ghq?", true) {
			_ = addUsingGhq(cmd, repoURL, rootPath)
			return
		}
		_ = addDirect(cmd, repoURL, rootPath)
	},
}

func addDirect(cmd *cobra.Command, repoURL, rootPath string) error {
	u, err := giturl.Parse(repoURL)
	if err != nil {
		printFatalln(cmd, err)
	}
	repoPath := filepath.Join(rootPath, strings.ReplaceAll(filepath.Join(u.Host, strings.TrimRight(u.Path, ".git")), "/", "__"))
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		printFatalln(cmd, err)
	}
	return cmdutil.Run("git", "clone", repoURL, repoPath)
}

func addUsingGhq(cmd *cobra.Command, repoURL, rootPath string) error {
	u, err := giturl.Parse(repoURL)
	if err != nil {
		printFatalln(cmd, err)
	}
	_ = cmdutil.Run("ghq", "get", repoURL)
	o, err := exec.Command("ghq", "list", "--full-path").Output()
	if err != nil {
		return err
	}
	repoPath := filepath.Join(rootPath, strings.ReplaceAll(filepath.Join(u.Host, strings.TrimRight(u.Path, ".git")), "/", "__"))
	var ghqPath string
	for _, path := range strings.Split(string(o), "\n") {
		if strings.Contains(path, filepath.Join(u.Host, strings.TrimRight(u.Path, ".git"))) {
			ghqPath = path
			break
		}
	}
	return cmdutil.Run("ln", "-s", ghqPath, repoPath)
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoAddCmd)
}
