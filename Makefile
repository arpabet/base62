VERSION := $(shell git describe --tags --always --dirty)
NOW := $(shell date +"%m-%d-%Y")

all: build

clean:
	go clean -i ./...

test:
	go test -cover ./...

build: test
	go build ./...
	go build -v -ldflags "-X main.Version=$(VERSION) -X main.Build=$(NOW)"  ./cmd/base62/...

update:
	go get -u ./...