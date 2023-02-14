DT := $(shell date +%Y%U)
REV := $(shell git rev-parse --short HEAD)
APP := $(shell basename $(CURDIR))
ARTIFACT := bin/$(APP)$(EXT)

TAGS ?= dev
GOFLAGS ?= -race -v
GOLDFLAGS ?= -X main.buildRevision=$(DT).$(REV)

.PHONY: all amd64 arm64 win mingw build linux release tidy updep

build: tidy
	go build $(GOFLAGS) -ldflags "$(GOLDFLAGS)" -tags="$(TAGS)" -o $(ARTIFACT) cmd/main.go

release:
	GOFLAGS="-trimpath" GOLDFLAGS="$(GOLDFLAGS) -s -w" TAGS="release" $(MAKE) build

linux:
	GOOS=linux $(MAKE) release

amd64:
	EXT=.x86-64 GOARCH=amd64 $(MAKE) linux

arm64:
	EXT=.aarch64 GOARCH=arm64 $(MAKE) linux

win:
	EXT=.exe CGO_ENABLED=1 CC=$(CURDIR)/mingw64/bin/x86_64-w64-mingw32-gcc.exe CXX=$(CURDIR)/mingw64/bin/x86_64-w64-mingw32-g++.exe $(MAKE) release

tidy: go.mod
	go mod tidy

updep: go.mod
	rm -f go.sum
	head -1 go.mod > /tmp/go.mod
	mv /tmp/go.mod go.mod
	go mod tidy

all: amd64 arm64 win
