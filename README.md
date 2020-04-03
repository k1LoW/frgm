# frgm [![Build Status](https://github.com/k1LoW/frgm/workflows/build/badge.svg)](https://github.com/k1LoW/frgm/actions)

`frgm` is a command snippets (fragments) manager.

## frgm export

### Export snippets as Alfred snippets

``` console
$ frgm export --to ~/Library/Application Support/Alfred/Alfred.alfredpreferences/snippets --format alfred
```

## frgm import

### Import Alfred snippets

``` console
$ frgm import --from ~/Library/Application Support/Alfred/Alfred.alfredpreferences/snippets --format alfred
```

## frgm list

``` console
$ frgm list
```

### zsh auto-complete from snippets using peco (Ctrl+j)

``` zsh
if exists frgm; then
    function peco-frgm() {
        BUFFER=$(frgm list | peco --query "$LBUFFER")
        CURSOR=$#BUFFER
        zle clear-screen
    }
    zle -N peco-frgm
    bindkey '^j' peco-frgm
fi
```

### zsh auto-complete from history and snippets using peco (Ctrl+r)

``` zsh
if exists peco; then
    function peco-select-history() {
        if exists frgm; then
            BUFFER=$((history -r -n 1 & frgm list) | \
                         peco --query "$LBUFFER")
        else
            BUFFER=$(history -r -n 1 | \
                         peco --query "$LBUFFER")
        fi
        CURSOR=$#BUFFER
        zle clear-screen
    }
    zle -N peco-select-history
    bindkey '^R' peco-select-history
fi
```
