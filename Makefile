
default: build

build:
	go build -o vacuum ./cmd

test:
	go test -v -race -timeout 30m ./...

