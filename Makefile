
default: build

build:
	go build -o vacuum ./cmd/vacuum

build-litter:
	go build -o litter ./cmd/litterbug

all: build build-litter

test:
	go test -v -race -timeout 30m ./...

