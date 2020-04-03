# frgm [![Build Status](https://github.com/k1LoW/frgm/workflows/build/badge.svg)](https://github.com/k1LoW/frgm/actions)

`frgm` is a command snippets (fragments) manager.

## frgm export

### Export snippets as [Alfred](https://www.alfredapp.com/) snippets

``` console
$ frgm export --to ~/Library/Application Support/Alfred/Alfred.alfredpreferences/snippets --format alfred
```

## frgm import

### Import [Alfred](https://www.alfredapp.com/) snippets

``` console
$ frgm import --from ~/Library/Application Support/Alfred/Alfred.alfredpreferences/snippets --format alfred
```

## frgm list

``` console
$ frgm list
```

### zsh auto-complete from snippets using [peco](https://github.com/peco/peco) (Ctrl+j)

``` zsh
function peco-select-snippets() {
    BUFFER=$(frgm list | peco --query "$LBUFFER")
    CURSOR=$#BUFFER
    zle clear-screen
}
zle -N peco-select-snippets
bindkey '^j' peco-select-snippets
```

### zsh auto-complete from history and snippets using [peco](https://github.com/peco/peco) (Ctrl+r)

``` zsh
function peco-select-history-and-snippets() {
    BUFFER=$((history -r -n 1 & frgm list) | peco --query "$LBUFFER")
    CURSOR=$#BUFFER
    zle clear-screen
}
zle -N peco-select-history-and-snippets
bindkey '^R' peco-select-history-and-snippets
```
