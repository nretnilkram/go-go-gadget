# go-go-gadget

Set of gadgets written in golang

## Install Package to different location

```
GOBIN=/usr/local/bin/ go install
```

### Setup

```
export GO_GO_GADGET_WORDS_FILE=$HOME/go-go-gadget/words/english_words.json
```

### Create Alias for go-go-gadget

```
alias ggg="~/go/bin/go-go-gadget"
```

### Help

```
go-go-gadget help
```

```
go-go-gadget is a set of command line tools.

Options:
  --help            Display go-go-gadget help
    -h, help
  --kubernetes      Take a string and k8s-ify it
    -k, --k8s, k8s
  --password        Return a password of desired length
    -p, --pw, pw, password
  --reverse         Take a string and reverse it
    -r, reverse
  --time            Display the current time
    -t, time
  --words           Return a set of words of desired length
    -w, words
```
