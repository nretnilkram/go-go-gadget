GOFMT_FILES?=$$(find . -name '*.go')

fmt:
	gofmt -w -l $(GOFMT_FILES)

init:
	brew install go@1.25

upgrade:
	go get -u ./...
	go mod tidy

install:
	GOBIN=$(HOME)/go/bin/ go install -v ./...

test:
	go test -v ./...

run:
	go run main.go

clean:
	rm -f go-go-gadget
