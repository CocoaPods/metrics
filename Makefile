CIRCLE_BUILD_NUM ?= dev
TAG?=0.0.$(CIRCLE_BUILD_NUM)-$(shell git rev-parse --short HEAD)

DATE = $(shell date "+%FT%T%z")

PREFIX?=$(shell pwd)
NAME := metrics
PKG := github.com/cocoapods/$(NAME)

BUILDTAGS ?= cgo

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := ${PREFIX}/cross

# Populate version variables
# Add to compile time flags
CTIMEVAR=-X $(PKG)/cmd.BuildDate=$(DATE) -X $(PKG)/cmd.Version=$(TAG)
GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"

# List the GOOS and GOARCH to build
GOOSARCHES = darwin/amd64 linux/arm64 linux/amd64

.PHONY: build
build: dist/$(NAME) ## Builds a dynamic executable or package

.PHONY: dist/$(NAME)
dist/$(NAME):
	@echo "+ $@"
	go build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $@ .

.PHONY: all
all: clean build lint test install

.PHONY: fmt
fmt:
	@echo "+ $@"
	@gometalinter --disable-all --enable gofmt ./...

.PHONY: test
test:
	@echo "+ $@"
	@gotestsum -- -tags "$(BUILDTAGS)" ./...

.PHONY: lint
lint:
	@echo "+ $@"
	@gometalinter ./...

.PHONY: cover
cover:
	@echo "+ $@"
	@go test -tags "$(BUILDTAGS) integration" -coverprofile=coverage.txt ./...

.PHONY: install
install:
	@echo "+ $@"
	go install -a -tags "$(BUILDTAGS)" ${GO_LDFLAGS} .

define buildpretty
mkdir -p $(BUILDDIR)/$(1)/$(2);
GOOS=$(1) GOARCH=$(2) go build \
	 -o $(BUILDDIR)/$(1)/$(2)/$(NAME) \
	 -a -tags "$(BUILDTAGS)" .;
endef

.PHONY: cross
cross: *.go
	@echo "+ $@"
	$(foreach GOOSARCH,$(GOOSARCHES), $(call buildpretty,$(subst /,,$(dir $(GOOSARCH))),$(notdir $(GOOSARCH))))

.PHONY: clean
clean:
	@echo "+ $@"
	$(RM) -r dist $(BUILDDIR)

