BIN_NAME ?= ./dist/neckless_$(shell uname -s | tr 'A-Z' 'a-z')_$(shell uname -m | sed 's/x86_64/amd64/'| sed 's/aarch64/arm64/')/neckless
VERSION ?= dev
GITCOMMIT ?= $(shell git rev-list -1 HEAD)
INSTALL_DIR ?= /usr/local/bin

all: test build

build: $(BIN_NAME) version

$(BIN_NAME): .goreleaser.yml
	goreleaser build --rm-dist

version: $(BIN_NAME)
	$(BIN_NAME) version	

install: $(BIN_NAME)
	cp $(BIN_NAME) $(INSTALL_DIR)

neckless:
	go build -ldflags "-s -w -X main.Version='$(VERSION)' -X main.GitCommit=$(GITCOMMIT)" -o $(BIN_NAME) github.com/mabels/neckless/cmd/neckless
	./neckless version

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


