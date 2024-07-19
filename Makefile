.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

APP=dobby
VERSION?=1.0
APP_EXECUTABLE="./out/$(APP)"
SRC_PACKAGES=$(shell go list ./... | grep -v "vendor")
SHELL=/bin/bash -o pipefail
GOPATH=$(shell go env GOPATH)
GOLANGCI_LINT_VERSION=v1.59.1
GOLANGCI_LINT=$(shell command -v golangci-lint 2> /dev/null)


BUILD?=$(shell git describe --always --dirty 2> /dev/null)
ifeq ($(BUILD),)
	BUILD=dev
endif

RICHGO=$(shell command -v richgo 2> /dev/null)
ifeq ($(RICHGO),)
	GO_BINARY=go
else
	GO_BINARY=richgo
endif

GO_TEST=$(shell command -v gotestsum 2> /dev/null)
ifeq ($(GO_TEST),)
	GO_TEST=$(GO_BINARY) test -mod=vendor $(SRC_PACKAGES) -coverprofile ./out/coverage -short -v
else
	GO_TEST=gotestsum --packages ${SRC_PACKAGES}
endif

ifdef CI_COMMIT_SHORT_SHA
	BUILD=$(CI_COMMIT_SHORT_SHA)
endif

setup-richgo:
ifeq ($(RICHGO),)
	$(GO_BINARY) get -u github.com/kyoh86/richgo
endif

setup-golangci-lint:
ifeq ($(GOLANGCI_LINT),)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin $(GOLANGCI_LINT_VERSION)
endif

SWAG=$(shell command -v swag 2> /dev/null)
setup-swag:
ifeq ($(SWAG),)
	$(GO_BINARY) install github.com/swaggo/swag/cmd/swag@latest
endif

setup: setup-golangci-lint setup-swag ensure-build-dir ## Setup environment

all: setup build

ensure-build-dir:
	mkdir -p out

ensure-vendor:
	go mod vendor

build-deps: ## Install dependencies
	go get
	go mod tidy
	go mod vendor

update-deps: ## Update dependencies
	go get -u

compile: ensure-build-dir ## Compile dobby
	$(GO_BINARY) build -ldflags "-X main.majorVersion=$(VERSION) -X main.minorVersion=${BUILD}" -o $(APP_EXECUTABLE) ./main.go

run: compile ## Run dobby
	./out/dobby server

compile-linux: ensure-build-dir ## Compile dobby for linux
	GOOS=linux GOARCH=amd64 $(GO_BINARY) build -ldflags "-X main.majorVersion=$(VERSION) -X main.minorVersion=${BUILD}" -o $(APP_EXECUTABLE) ./main.go

build: fmt lint test compile ## Build the application

compress: compile ## Compress the binary
	upx $(APP_EXECUTABLE)

fmt:
	$(GO_BINARY) fmt $(SRC_PACKAGES)

lint: setup-golangci-lint
	$(GOLANGCI_LINT) run -v

test: ensure-build-dir ## Run tests
	ENVIRONMENT=test $(GO_TEST)

test-cover-html: ## Run tests with coverage
	mkdir -p ./out
	@echo "mode: count" > coverage-all.out
	$(foreach pkg, $(SRC_PACKAGES),\
	ENVIRONMENT=test $(GO_BINARY) test -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	$(GO_BINARY) tool cover -html=coverage-all.out -o out/coverage.html

swagger-docs: setup-swag ## Generate swagger docs
	$(SWAG) init

dockerfile-security: ## Dockerfile OPA
	docker run --rm -v $(PWD):/project openpolicyagent/conftest test --policy dockerfile-security.rego Dockerfile