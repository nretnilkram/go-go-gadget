# go-go-gadget

Set of gadgets written in golang

## Install Package to different location

```sh
GOBIN=$HOME/go/bin/ go install
```

### Create Alias for go-go-gadget

```sh
alias ggg="~/go/bin/go-go-gadget"
```

### Execution

You can run the app by simply typing `go-go-gadget` and following the prompts.  Or you can pass all the parameters at once.

```sh
go-go-gadget
go-go-gadget password 25
go-go-gadget -w 1000
```

```sh
go run main.go
```

### Help

```sh
go-go-gadget -h
go-go-gadget --help
```

### Testing

```sh
GOBIN=$HOME/go/bin/ go install

PATH="$PATH:$HOME/go/bin" go test -v ./...

PATH="$PATH:$HOME/go/bin" go test -v ./pkg/...

PATH="$PATH:$HOME/go/bin" go test -v main_test.go
```
