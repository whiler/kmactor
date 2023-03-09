DT := $(shell date +%Y%U)
REV := $(shell git rev-parse --short HEAD)
APP := $(shell basename $(CURDIR))
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
EXT :=
ifeq ($(GOOS), windows)
	EXT := .exe
endif
ARTIFACT := bin/$(APP)-$(GOOS)-$(GOARCH)$(EXT)

TAGS ?= dev
GOFLAGS ?= -race -v
GOLDFLAGS ?= -X main.buildRevision=$(DT).$(REV)

.PHONY: all amd64 arm64 build release tidy updep

build:
	CGO_ENABLED=1 go build $(GOFLAGS) -ldflags "$(GOLDFLAGS)" -tags="$(TAGS)" -o $(ARTIFACT) cmd/main.go

release:
	GOFLAGS="-trimpath" GOLDFLAGS="$(GOLDFLAGS) -s -w" TAGS="release" $(MAKE) build

amd64:
	GOARCH=amd64 $(MAKE) release

arm64:
	GOARCH=arm64 $(MAKE) release

tidy: go.mod
	go mod tidy

updep: go.mod
	rm -f go.sum
	head -1 go.mod > /tmp/go.mod
	mv /tmp/go.mod go.mod
	go mod tidy

all: amd64 arm64
