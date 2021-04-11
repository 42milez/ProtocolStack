include .env

PROJECTNAME = $(shell basename "$(PWD)")

GOBIN = "$(PWD)/bin"

STDERR = /tmp/.$(PROJECTNAME)-stderr.txt
MAKEFLAGS += --silent

TCP_CLIENT_FILES = src/tcp_client.go
TCP_CLIENT_BIN = tcp_client

TCP_SERVER_FILES = src/tcp_server.go
TCP_SERVER_BIN = tcp_server

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

#  Make Commands
# --------------------------------------------------
## build: Build the program.
.PHONY: build
build: go-build

## clean: Clean up caches.
.PHONY: clean
clean: go-clean

## compile: Clean up caches, resolve dependencies, build the program.
.PHONY: compile
compile:
	@-rm $(STDERR)
	@-touch $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## resolve: Resolve dependencies.
.PHONY: resolve
resolve: go-mod

#  Go Commands
# --------------------------------------------------
.PHONY: go-build
go-build:
	@echo "  >  Building binary..."
	@mkdir -p $(GOBIN)
	@go build -o $(GOBIN)/$(TCP_CLIENT_BIN) $(TCP_CLIENT_FILES)
	@go build -o $(GOBIN)/$(TCP_SERVER_BIN) $(TCP_SERVER_FILES)

.PHONY: go-clean
go-clean:
	@echo "  >  Cleaning build cache"
	@go clean

.PHONY: go-compile
go-compile: go-clean go-mod go-build

.PHONY: go-mod
go-mod:
	@echo "  >  Checking if there is any missing dependencies..."
	@go mod tidy
