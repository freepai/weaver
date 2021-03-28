# Makefile for building CoreDNS
GITCOMMIT:=$(shell git describe --dirty --always)
BINARY:=weaver
SYSTEM:=
CHECKS:=check
BUILDOPTS:=-v
GOPATH?=$(HOME)/go
MAKEPWD:=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))
CGO_ENABLED:=0

.PHONY: all
all: build

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -o bin/$(BINARY)

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -c -covermode=count -coverpkg ./...

.PHONY: run
run: build
	./bin/${BINARY}

.PHONY: clean
clean:
	go clean
	rm -f bin