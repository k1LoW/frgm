package history

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/k1LoW/frgm/snippet"
)

type History struct{}

// New return *History
func New() *History {
	return &History{}
}

var trimRe = regexp.MustCompile(`^\s*[0-9*]+\s+`)

func (h *History) Decode(stdin io.Reader, group string) (snippet.Snippets, error) {
	snippets := snippet.Snippets{}
	in := bufio.NewReader(stdin)
	for {
		line, err := in.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return snippets, err
		}
		command := strings.TrimRight(trimRe.ReplaceAllString(line, ""), "\n")
		if command == "" {
			continue
		}
		name := command
		if strings.Count(command, "#") == 1 {
			splited := strings.Split(command, "#")
			command = strings.TrimSpace(splited[0])
			name = strings.TrimSpace(splited[1])
		}
		uid, err := genUID(line)
		if err != nil {
			return snippets, err
		}
		snippets = append(snippets, snippet.New(uid, group, name, command, "", "", []string{}))
	}
	return snippets, nil
}

func genUID(line string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, line); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	return fmt.Sprintf("history-%s", s[:12]), nil
}
