include .env

PROJECT_NAME = $(shell basename "$(PWD)")

GOBIN = "$(PWD)/bin"

STDERR = /tmp/$(PROJECT_NAME)-stderr.txt
MAKEFLAGS += --silent

TCP_CLIENT_FILES = src/tcp_client.go
TCP_CLIENT_BIN = tcp_client

TCP_SERVER_FILES = src/tcp_server.go
TCP_SERVER_BIN = tcp_server

# https://golang.org/cmd/compile/#hdr-Command_Line
# https://golang.org/doc/gdb#Introduction
ifeq ($(RELEASE), true)
  BUILD_FLAGS = ""
else
  BUILD_FLAGS = -gcflags=all="-N -l"
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

## compile: clean up caches, resolve dependencies, and build the program
.PHONY: compile
compile:
	@-rm -f $(STDERR)
	@-touch $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)

## resolve: resolve dependencies
.PHONY: resolve
resolve: go-mod

## test: run all tests
.PHONY: test
test:
	@go test -v ./src/network

#  Go Commands
# --------------------------------------------------
.PHONY: go-build
go-build:
	@echo "Building binary..."
	@mkdir -p $(GOBIN)
	@go build $(BUILD_FLAGS) -o $(GOBIN)/$(TCP_CLIENT_BIN) $(TCP_CLIENT_FILES)
	@go build $(BUILD_FLAGS) -o $(GOBIN)/$(TCP_SERVER_BIN) $(TCP_SERVER_FILES)

.PHONY: go-clean
go-clean:
	@echo "Cleaning build cache..."
	@go clean

.PHONY: go-compile
go-compile: go-clean go-mod go-build

.PHONY: go-mod
go-mod:
	@echo "Checking if there is any missing dependencies..."
	@go mod tidy
