# go-go-gadget

Set of cli gadgets written in golang

Requires **Go 1.27+**.

## Install Package

From a clone of this repo:

```sh
GOBIN=$HOME/go/bin/ go install ./...
```

Or without cloning:

```sh
GOBIN=$HOME/go/bin/ go install github.com/nretnilkram/go-go-gadget@latest
```

### Create Alias for go-go-gadget

```sh
alias go-go-gadget="~/go/bin/go-go-gadget"
alias ggg="go-go-gadget"
```

### Setup Autocomplete

```sh
go-go-gadget completion zsh > /path/to/completion/dir/go_go_gadget_auto_complete
```

### Execution

Pass parameters on the command line (most commands are non-interactive):

```sh
go-go-gadget [help|--help|-h]

# Password generation
go-go-gadget [password|pw|p] [--length|-l] 25                  # symbols included by default
go-go-gadget [password|pw|p] [--length|-l] 25 --symbols=false  # disable symbols
go-go-gadget p -s=false -l=64

# String manipulation
go-go-gadget k8s "Nretnil Kram"
go-go-gadget reverse "hello world"
go-go-gadget symsub "password"
go-go-gadget inspect "test123" [--digits|-d]

# Words and dates
go-go-gadget words [--count|-c] 1000
go-go-gadget words --count 8 --adjective-weight 1 --animal-weight 3 --verb-weight 1
go-go-gadget now [--time|-t] [--colon|-c|--dots|-d|--slashes|-s|--unix|-u|--raw|-r]

# Grit - Git repository management
go-go-gadget grit init
go-go-gadget grit add-repo <path>
go-go-gadget grit add-all-repos
go-go-gadget grit remove-repo <name-or-pattern>   # glob; use % as wildcard (e.g. gg-phoenix-%)
go-go-gadget grit remove-repo 'gg-phoenix-*'      # quoted * also works
go-go-gadget grit config
go-go-gadget grit history
go-go-gadget grit reset
go-go-gadget grit destroy
go-go-gadget grit pull
go-go-gadget grit [--help|-h]

# Group - run any command in parallel across grit config directories
go-go-gadget group ls -la
go-go-gadget group npm install
go-go-gadget group make test

# Limit concurrent grit/group actions (default: min(GOMAXPROCS*2, 16))
# Also configurable via max_concurrent in .grit/config.yml
export GRIT_MAX_CONCURRENT=10

# Utilities
go-go-gadget [utilities|u|util|utils] quick-branch [name]
go-go-gadget [utilities|u|util|utils] quick-commit
go-go-gadget [utilities|u|util|utils] empty-commit
go-go-gadget [utilities|u|util|utils] tf-list-resources <file1> [file2...]
go-go-gadget [utilities|u|util|utils] semver "1.2.0"
go-go-gadget [utilities|u|util|utils] bump-chart major|minor|patch
go-go-gadget utilities tf-list-resources test.tf
go-go-gadget utilities semver "1.2.0"
go-go-gadget utilities semver "A.2.0"

# Kubernetes (m8s)
go-go-gadget m8s deployment [--alpine|-a|--busybox|-b|--ubuntu|-u] [--namespace|-n default] [--image|-i <image>]
go-go-gadget m8s pod -u -n shared
go-go-gadget m8s connection -i debian
go-go-gadget m8s connection -i debian --command /bin/bash
```

## Run Locally

```sh
go run main.go
```

Useful Makefile targets:

```sh
make help           # list targets
make test           # go test -v ./...
make test-coverage  # coverage profile
make lint           # golangci-lint
make install        # install into $HOME/go/bin
make gamut          # fmt, lint, test, install, etc.
```

## Update Dependencies

```sh
go get -u ./...
# or: make upgrade
```

## Help

```sh
go-go-gadget -h
go-go-gadget --help
```

```sh
Go Go Gadget...Help!

* * * * * * * * * * * * *
 
Go Go Gadget is a set of helpful CLI tools

You can use go-go-gadget to run useful tools from the terminal.

Usage:
  go-go-gadget [flags]
  go-go-gadget [command]

Available Commands:
  completion             Generate the autocompletion script for the specified shell
  generate-documentation Generate Documentation
  grit                   Run git commands on multiple repositories
  group                  Run a command on multiple directories (uses grit config to determine directories)
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
  -h, --help      help for go-go-gadget
  -v, --version   version for go-go-gadget

Use "go-go-gadget [command] --help" for more information about a command.
```

### Generate Documentation

```sh
go-go-gadget [generate-documentation|docs|documentation] [md|rest|yaml]
```

## Testing

```sh
go test -v ./...

go test -v ./pkg/...

# or: make test
```
