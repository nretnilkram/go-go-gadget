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

### Setup Autocomplete

```sh
go-go-gadget completion zsh > (/path/to/completion/dir)/go_go_gadget_auto_complete
```

### Execution

You can run the app by simply typing `go-go-gadget` and following the prompts.  Or you can pass all the parameters at once.

```sh
go-go-gadget [help|--help|-h]

go-go-gadget [password|pw] [--length|-l] 25
go-go-gadget words [--count|-c] 1000
go-go-gadget k8s "Nretnil Kram"

go-go-gadget grit init
go-go-gadget grit add-all
go-go-gadget grit reset
go-go-gadget grit config
go-go-gadget grit pull
go-go-gadget grit [--help|-h]

go-go-gadget [utilities|u] quick-branch
go-go-gadget [utilities|u] quick-commit
go-go-gadget [utilities|u] tf-list-resources
go-go-gadget [utilities|u] semver

go-go-gadget m8s deployment
go-go-gadget m8s pod -u -n shared
go-go-gadget m8s connection -i debian
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
  now         Show todays date
  password    Generate a password
  reverse     Reverses a string
  symsub      Substitute symbols into a string
  utilities   Useful utility commands
  words       Create list of words

Flags:
  -h, --help      help for go_go_gadget
  -v, --version   version for go_go_gadget

Use "go_go_gadget [command] --help" for more information about a command.
```

### Generate Documentation

```sh
go-go-gadget docs [md|rest|yaml]
```

## Testing

```sh
go test -v ./...

go test -v ./pkg/...

go test -v main_test.go
```
