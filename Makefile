
default: build

build:
	env GO111MODULE=on go build -o vacuum ./cmd/vacuum

build-litter:
	env GO111MODULE=on go build -o litter ./cmd/litterbug

all: build build-litter

test:
	go test -v -race -timeout 30m ./...

