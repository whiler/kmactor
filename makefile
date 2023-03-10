DT := $(shell date +%Y%U)
REV := $(shell git rev-parse --short HEAD)
APP := $(shell basename $(CURDIR))
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
EXT :=
ifeq ($(GOOS),windows)
	EXT := .exe
endif
ARTIFACT := bin/$(APP)-$(GOOS)-$(GOARCH)$(EXT)
AMD64CC :=
ARM64CC :=

TAGS ?= dev
GOFLAGS ?= -race -v
GOLDFLAGS ?= -X main.buildRevision=$(DT).$(REV)

ifeq ($(shell go env GOHOSTOS), windows)
	AMD64CC = x86_64-w64-mingw32-gcc
	ARM64CC = aarch64-w64-mingw32-gcc
else ifeq ($(shell go env GOHOSTOS), linux)
ifeq ($(shell go env GOHOSTARCH), amd64)
	ARM64CC = aarch64-linux-gnu-gcc
else
	AMD64CC = x86_64-linux-gnu-gcc
endif
endif

.PHONY: all amd64 arm64 build release tidy updep

build:
	CGO_ENABLED=1 go build $(GOFLAGS) -ldflags "$(GOLDFLAGS)" -tags="$(TAGS)" -o $(ARTIFACT) cmd/*.go

release:
	GOFLAGS="-trimpath" GOLDFLAGS="$(GOLDFLAGS) -s -w" TAGS="release" $(MAKE) build

amd64:
	GOARCH=amd64 CC=$(AMD64CC) $(MAKE) release

arm64:
	GOARCH=arm64 CC=$(ARM64CC) $(MAKE) release

tidy: go.mod
	go mod tidy

updep: go.mod
	rm -f go.sum
	head -1 go.mod > /tmp/go.mod
	mv /tmp/go.mod go.mod
	go mod tidy

all: amd64 arm64
