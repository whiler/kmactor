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
	GOLDFLAGS = -H=windowsgui -X main.buildRevision=$(DT).$(REV)
else ifeq ($(shell go env GOHOSTOS), linux)
ifeq ($(shell go env GOHOSTARCH), amd64)
	ARM64CC = aarch64-linux-gnu-gcc
else
	AMD64CC = x86_64-linux-gnu-gcc
endif
endif

.PHONY: all amd64 arm64 build macapp release tidy updep

build:
	CGO_ENABLED=1 go build $(GOFLAGS) -ldflags "$(GOLDFLAGS)" -tags="$(TAGS)" -o $(ARTIFACT) cmd/*.go

release:
	GOFLAGS="-trimpath" GOLDFLAGS="$(GOLDFLAGS) -s -w" TAGS="release" $(MAKE) build

amd64:
ifeq ($(shell go env GOHOSTOS), windows)
	goversioninfo -64
endif
	GOARCH=amd64 CC=$(AMD64CC) $(MAKE) release
ifeq ($(shell go env GOHOSTOS), windows)
	rm -f resource.syso
endif

arm64:
ifeq ($(shell go env GOHOSTOS), windows)
	goversioninfo -64 -arm
endif
	GOARCH=arm64 CC=$(ARM64CC) $(MAKE) release
ifeq ($(shell go env GOHOSTOS), windows)
	rm -f resource.syso
endif

tidy: go.mod
	go mod tidy

updep: go.mod
	rm -f go.sum
	head -1 go.mod > /tmp/go.mod
	mv /tmp/go.mod go.mod
	go mod tidy

macapp: kmactor.app
	zip bin/$(APP).zip -r kmactor.app

all: amd64 arm64
