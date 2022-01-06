BIN_NAME ?= "./neckless"
VERSION ?= dev
GITCOMMIT ?= $(shell git rev-list -1 HEAD)

all: test build

build:
	goreleaser build --rm-dist
	$(BIN_NAME) version
	dist/neckless_linux_amd64/neckless
	# go build -ldflags "-s -w -X main.Version='$(VERSION)' -X main.GitCommit=$(GITCOMMIT)" -o $(BIN_NAME) github.com/mabels/neckless/cmd/neckless

test:
	go test github.com/mabels/neckless/key
	go test github.com/mabels/neckless/symmetric
	go test github.com/mabels/neckless/asymmetric
	go test github.com/mabels/neckless/casket
	go test github.com/mabels/neckless/member
	go test github.com/mabels/neckless/pearl
	go test github.com/mabels/neckless/kvpearl
	go test github.com/mabels/neckless/gem
	go test github.com/mabels/neckless/necklace
	go test github.com/mabels/neckless/cmd/neckless


