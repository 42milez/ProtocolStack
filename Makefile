include .env

PROJECT_NAME := $(shell basename "$(PWD)")

GOBIN := ./bin

STDERR := /tmp/$(PROJECT_NAME)-stderr.txt
MAKEFLAGS += --silent

TCP_SERVER_FILES := src/tcp_server.go
TCP_SERVER_BIN := tcp_server

# https://golang.org/cmd/compile/#hdr-Command_Line
# https://golang.org/doc/gdb#Introduction
ifeq ($(RELEASE), true)
  BUILD_FLAGS :=
else
  BUILD_FLAGS := -gcflags=all="-N -l"
endif

.PHONY: help
help: Makefile
	@echo "Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n "s/^##//p" $< | column -t -s ":" |  sed -e "s/^/ /"
	@echo

#  Make Commands
# --------------------------------------------------
## build: build project
.PHONY: build
build: go-build

## clean: clean up caches
.PHONY: clean
clean: go-clean

## compile: clean up caches, resolve dependencies, and build the application
.PHONY: compile
compile:
	@-rm -f $(STDERR)
	@-touch $(STDERR)
	@$(MAKE) -s go-compile 2> $(STDERR)

## fmt: run formatter
.PHONY: fmt
fmt: go-fmt

## lint: run linters (golangci-lint)
.PHONY: lint
lint: go-lint

## resolve: resolve dependencies
.PHONY: resolve
resolve: go-mod

## test: run all tests
.PHONY: test
test:
	@go test -race -covermode=atomic -coverprofile=coverage.out -v ./src/network

#  Go Commands
# --------------------------------------------------
.PHONY: go-build
go-build:
	@echo "> building binary..."
	@mkdir -p $(GOBIN)
	@go build $(BUILD_FLAGS) -o $(GOBIN)/$(TCP_SERVER_BIN) $(TCP_SERVER_FILES)

.PHONY: go-clean
go-clean:
	@echo "> cleaning up build caches..."
	@go clean

.PHONY: go-compile
go-compile: go-clean go-mod go-vet go-build

.PHONY: go-fmt
go-fmt:
	@echo "> reformatting code..."
	@go fmt $(dir $(abspath $(firstword $(MAKEFILE_LIST))))src/...

.PHONY: go-lint
go-lint:
	@echo "> running linters..."
	@golangci-lint run

.PHONY: go-mod
go-mod:
	@echo "> checking if there is any missing dependencies..."
	@go mod tidy

.PHONY: go-vet
go-vet:
	@echo "> vetting source code..."
	@go vet $(dir $(abspath $(firstword $(MAKEFILE_LIST))))src/...
