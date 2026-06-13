VERSION := $(shell git describe --tags --always --dirty)
NOW := $(shell date +"%m-%d-%Y")

all: build

clean:
	go clean -i ./...

vet:
	go vet ./...

test: vet
	go test -cover ./...

fuzz:
	go test -run='^$$' -fuzz='FuzzRoundTrip' -fuzztime=30s .
	go test -run='^$$' -fuzz='FuzzDecode$$' -fuzztime=30s .

build: test
	go build ./...
	go build -v -ldflags "-X main.Version=$(VERSION) -X main.Build=$(NOW)"  ./cmd/base62/...

update:
	go get -u ./...