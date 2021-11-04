
default: build

build:
	env GO111MODULE=on go build -o vacuum ./cmd

test:
	go test -v -race -timeout 30m ./...

