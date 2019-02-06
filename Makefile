.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

APP=dobby
SRC_PACKAGES=$(shell go list ./... | grep -v "vendor")
VERSION?=1.0
BUILD?=$(shell git describe --always --dirty 2> /dev/null)
DEP:=$(shell command -v dep 2> /dev/null)
GOLINT:=$(shell command -v golint 2> /dev/null)
APP_EXECUTABLE="./out/$(APP)"
RICHGO=$(shell command -v richgo 2> /dev/null)
GOMETA_LINT=$(shell command -v golangci-lint 2> /dev/null)
GOLANGCI_LINT_VERSION=v1.12.5
GO111MODULE=off
SHELL=/bin/bash -o pipefail

ifeq ($(GOMETA_LINT),)
	GOMETA_LINT=$(shell command -v $(PWD)/bin/golangci-lint 2> /dev/null)
endif

ifeq ($(RICHGO),)
	GO_BINARY=go
else
	GO_BINARY=richgo
endif

ifeq ($(BUILD),)
	BUILD=dev
endif

ifdef CI_COMMIT_SHORT_SHA
	BUILD=$(CI_COMMIT_SHORT_SHA)
endif

all: setup build

ensure-build-dir:
	mkdir -p out

build-deps: ## Install dependencies
	dep ensure -v

compile: ensure-build-dir ## Compile dobby
	$(GO_BINARY) build -ldflags "-X main.majorVersion=$(VERSION) -X main.minorVersion=${BUILD}" -o $(APP_EXECUTABLE) ./main.go

run: compile ## Run dobby
	./out/dobby server

compile-linux: ensure-build-dir ## Compile dobby for linux
	GOOS=linux GOARCH=amd64 $(GO_BINARY) build -ldflags "-X main.majorVersion=$(VERSION) -X main.minorVersion=${BUILD}" -o $(APP_EXECUTABLE) ./main.go

build: build-deps fmt vet lint-all compile ## Build the application

compress: compile ## Compress the binary
	upx $(APP_EXECUTABLE)

fmt:
	$(GO_BINARY) fmt $(SRC_PACKAGES)

vet:
	$(GO_BINARY) vet $(SRC_PACKAGES)

setup-golangci-lint:
ifeq ($(GOMETA_LINT),)
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(GOLANGCI_LINT_VERSION)
endif

setup: setup-golangci-lint ensure-build-dir ## Setup environment
ifeq ($(DEP),)
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
ifeq ($(GOLINT),)
	$(GO_BINARY) get -u golang.org/x/lint/golint
endif
ifeq ($(RICHGO),)
	$(GO_BINARY) get -u github.com/kyoh86/richgo
endif

lint-all: lint setup-golangci-lint
	$(GOMETA_LINT) run

lint:
	./scripts/lint $(SRC_PACKAGES)

test-all: test test.integration

test: ensure-build-dir ## Run tests
	ENVIRONMENT=test $(GO_BINARY) test $(SRC_PACKAGES) -p=1 -coverprofile ./out/coverage -short -v | grep -vi "start" | grep -vi "no test files"

test-cover-html: ## Run tests with coverage
	mkdir -p ./out
	@echo "mode: count" > coverage-all.out
	$(foreach pkg, $(SRC_PACKAGES),\
	ENVIRONMENT=test $(GO_BINARY) test -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	$(GO_BINARY) tool cover -html=coverage-all.out -o out/coverage.html