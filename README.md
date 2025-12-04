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

# Password generation
go-go-gadget [password|pw|p] [--length|-l] 25 [--symbols|-s]

# String manipulation
go-go-gadget k8s "Nretnil Kram"
go-go-gadget reverse "hello world"
go-go-gadget symsub "password"
go-go-gadget inspect "test123" [--digits|-d]

# Words and dates
go-go-gadget words [--count|-c] 1000
go-go-gadget now [--time|-t] [--colon|-c|--dots|-d|--slashes|-s|--unix|-u|--raw|-r]

# Grit - Git repository management
go-go-gadget grit init
go-go-gadget grit add-repo <path>
go-go-gadget grit add-all-repos
go-go-gadget grit remove-repo <name>
go-go-gadget grit config
go-go-gadget grit history
go-go-gadget grit reset
go-go-gadget grit destroy
go-go-gadget grit pull
go-go-gadget grit [--help|-h] [--synchronous|-s]

# Utilities
go-go-gadget [utilities|u|util|utils] quick-branch [name]
go-go-gadget [utilities|u|util|utils] quick-commit
go-go-gadget [utilities|u|util|utils] empty-commit
go-go-gadget [utilities|u|util|utils] tf-list-resources <file1> [file2...]
go-go-gadget [utilities|u|util|utils] semver "1.2.0"

# Kubernetes (m8s)
go-go-gadget m8s deployment [--alpine|-a|--busybox|-b|--ubuntu|-u] [--namespace|-n default] [--image|-i <image>]
go-go-gadget m8s pod -u -n shared
go-go-gadget m8s connection -i debian
```

## Run Locally

```sh
go run main.go
```

## Update Dependencies

```sh
go get -u ./...
```

## Help

```sh
go-go-gadget -h
go-go-gadget --help
```

```sh
Go Go Gadget is a set of helpful CLI tools

You can use go_go_gadget to run useful tools from the terminal.

Usage:
  go_go_gadget [flags]
  go_go_gadget [command]

Available Commands:
  completion             Generate the autocompletion script for the specified shell
  generate-documentation Generate Documentation
  grit                   Run git commands on multiple repositories
  help                   Help about any command
  inspect                Inspects a string
  k8s                    Kubernetesify a string
  m8s                    Useful k8s commands
  now                    Print todays date
  password               Generate a password
  reverse                Reverses a string
  symsub                 Substitute symbols
  utilities              Useful utility commands
  words                  Create list of words

Flags:
  -h, --help      help for go_go_gadget
  -v, --version   version for go_go_gadget

Use "go_go_gadget [command] --help" for more information about a command.
```

### Generate Documentation

```sh
go-go-gadget [generate-documentation|docs|documentation] [md|rest|yaml]
```

## Testing

```sh
go test -v ./...

go test -v ./pkg/...

go test -v main_test.go
```
