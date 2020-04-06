# frgm [![Build Status](https://github.com/k1LoW/frgm/workflows/build/badge.svg)](https://github.com/k1LoW/frgm/actions)

`frgm` is a command snippets (fragments) manager.

## Usage

### frgm init

Initialize frgm.

- Create or update config.toml
- Set and create `global.snippets_path` ( create file or make directory )

``` console
$ frgm init
```

### frgm export

#### Export snippets as [Alfred](https://www.alfredapp.com/) snippets

``` console
$ frgm export --to ~/Library/Application Support/Alfred/Alfred.alfredpreferences/snippets --format alfred
```

#### Export snippets as Markdown document.

``` console
$ frgm export --to /path/to/snippets.md --format md
```

### frgm import

#### Import [Alfred](https://www.alfredapp.com/) snippets

``` console
$ frgm import --from ~/Library/Application Support/Alfred/Alfred.alfredpreferences/snippets --format alfred
```

### frgm list

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
