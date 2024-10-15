# go-go-gadget

Set of cli gadgets written in golang

## Install Package

```sh
GOBIN=$HOME/go/bin/ go install
```

### Create Alias for go-go-gadget

```sh
alias go-go-gadget="~/go/bin/go-go-gadget"
alias ggg="go-go-gadget"
```

### Execution

You can run the app by simply typing `go-go-gadget` and following the prompts.  Or you can pass all the parameters at once.

```sh
go-go-gadget
go-go-gadget password --length 25
go-go-gadget words --count 1000
go-go-gadget k8s "Nretnil Kram"

go-go-gadget grit init
go-go-gadget grit add-all
go-go-gadget grit reset
go-go-gadget grit config
go-go-gadget grit pull

go-go-gadget utilities quick-branch
go-go-gadget utilities quick-commit
```

## Run Locally

```sh
go run main.go
```

## Help

```sh
go-go-gadget -h
go-go-gadget --help
```

```sh
Go Go Gadget is a set of helpful CLI tools

One can use go_go_gadget to run useful tools from the terminal

Usage:
  go_go_gadget [flags]
  go_go_gadget [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  grit        Run git commands on multiple repositories at once
  help        Help about any command
  inspect     Inspects a string
  k8s         Kubernetesify a string
  password    Generate a password
  reverse     Reverses a string
  symsub      Substitute symbols into a string
  words       create list of words

Flags:
  -h, --help      help for go_go_gadget
  -v, --version   version for go_go_gadget

Use "go_go_gadget [command] --help" for more information about a command.
```

## Testing

```sh
PATH="$PATH:$HOME/go/bin" go test -v ./...

PATH="$PATH:$HOME/go/bin" go test -v ./pkg/...

PATH="$PATH:$HOME/go/bin" go test -v main_test.go
```
