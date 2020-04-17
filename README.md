# frgm [![Build Status](https://github.com/k1LoW/frgm/workflows/build/badge.svg)](https://github.com/k1LoW/frgm/actions)

`frgm` is a meta snippet (fragment) manager.

Key features of `frgm` are:

- **[Can export snippets in other snippet tool format](#export-frgm-snippets-as--snippets--frgm-export-)**.
- **[Can document snippets](#export-snippets-as-markdown-document)**.
- **[Can use as command-line snippets](##use-frgm-snippets--frgm-list-)**.

## Usage

### Initialize frgm ( `frgm init` )

Initialize frgm.

- Create or update config.toml
- Set and create `global.snippets_path` ( create file or make directory )

``` console
$ frgm init
```

### Write snippets ( `frgm edit` )

Write snippets.

The format of the frgm snippet is the following,

``` yaml
---
# Default group name of snippets
# Default is snippets file name
group: my-group
snippets:
  # Unique identifier of snippet
  # Default is automatically generated
- uid: frgm-1ca779b751a5
  # Group name of snippet
  # Default is default group name of snippets
  group: command
  # Name of snippet
  name: Delete branch already merged
  # Content (command) of snippet
  content: git branch --merged | grep -v master | xargs git branch -d
  # Description of snippet
  desc: |
    1. lists the merged branches
    2. delete all merged branches except the master branch
    ref: https://example.com/path/to/link
  # Labels
  labels:
  - git
  - cleanup
- name: ping
  content: ping 8.8.8.8
[...]
```

You can use the `frgm edit` command to edit snippets of `global.snippets_path` using the editor specified in $EDITOR.

``` console
$ EDITOR=emacs frgm edit
```

### Export frgm snippets as * snippets ( `frgm export` )

#### Export snippets as [Alfred](https://www.alfredapp.com/) snippets

``` console
$ frgm export --to ~/Library/Application Support/Alfred/Alfred.alfredpreferences/snippets --format alfred
```

#### Export snippets as [pet](https://github.com/knqyf263/pet) snippets

``` console
$ frgm export --to /path/to/pet.toml --format pet
```

#### Export snippets as Markdown document.

``` console
$ frgm export --to /path/to/snippets.md --format md
```

**Key Mapping:**

| frgm snippet key | Default / Required | [Alfred](https://www.alfredapp.com/) | [pet](https://github.com/knqyf263/pet) |
| --- | --- | --- | --- |
| `uid:` | Default is automatically generated | `uid` | - |
| `group:` | Default is default group name of snippets or file name | directory | - |
| `name:` | required | `name` | `description:` |
| `desc:` | | - | - |
| `content:` | required | `snippet` | `command:` |
| `output:` | | - | `output:` |
| `labels:` | | `keyword` | `tag:` |

#### Fill uid, group ( `frgm fill` )

Fill and freeze `uid:` and `group:`.

``` console
$ frgm fill
```

**before:**

``` yaml
# my-group.yml
snippets:
- name: ping
  content: ping 8.8.8.8
```

**after:**

``` yaml
# my-group.yml
snippets:
- uid: frgm-6aa9d75f9d83
  group: my-group
  name: ping
  content: ping 8.8.8.8
```

### Import * snippets ( `frgm import` )

#### Import [Alfred](https://www.alfredapp.com/) snippets

``` console
$ frgm import --from ~/Library/Application Support/Alfred/Alfred.alfredpreferences/snippets --format alfred
```

### Convert * snippets ( `frgm convert` )

#### Convert `history` output into frgm snippets

``` console
$ history | tail -1 | ./frgm convert --format history --group from-history
snippets:
- uid: history-8cc0c8477ec0
  group: from-history
  name: git branch --merged | grep -v master | xargs git branch -d
  content: git branch --merged | grep -v master | xargs git branch -d
```

#### Convert [Alfred](https://www.alfredapp.com/) snippet file into frgm snippet

``` console
$ cat /path/to/alfredsnippet.json | frgm convert --format alfred --group search
snippets:
- uid: frgm-46c29e119523
  group: search
  name: Search log file
  content: lsof -c nginx | grep -v .so | grep .log | awk '{print $9}' | sort | uniq
  labels:
  - log
```

#### Convert [pet](https://github.com/knqyf263/pet) snippets file into frgm snippets

``` console
$ cat /path/to/pet-snippets.toml. | frgm convert --format pet --group pet
snippets:
- uid: pet-df7bb29f9681
  group: pet
  name: ping
  content: ping 8.8.8.8
  labels:
  - network
  - google
- uid: pet-584a331fd6b0
  group: pet
  name: hello
  content: echo hello
  output: hello
  labels:
  - sample
```

### Use frgm snippets ( `frgm list` )

``` console
$ frgm list
```

#### zsh auto-complete from snippets using [peco](https://github.com/peco/peco) (Ctrl+j)

``` zsh
function peco-select-snippets() {
    BUFFER=$(frgm list | peco --query "$LBUFFER")
    CURSOR=$#BUFFER
    zle clear-screen
}
zle -N peco-select-snippets
bindkey '^j' peco-select-snippets
```

#### zsh auto-complete from history and snippets using [peco](https://github.com/peco/peco) (Ctrl+r)

``` zsh
function peco-select-history-and-snippets() {
    BUFFER=$((history -r -n 1 & frgm list) | peco --query "$LBUFFER")
    CURSOR=$#BUFFER
    zle clear-screen
}
zle -N peco-select-history-and-snippets
bindkey '^R' peco-select-history-and-snippets
```

## Install

**homebrew tap:**

```console
$ brew install k1LoW/tap/frgm
```

**manually:**

Download binany from [releases page](https://github.com/k1LoW/frgm/releases)

**go get:**

```console
$ go get github.com/k1LoW/frgm
```
