NAME=caretaker

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := bin

# VERSION defines the project version for the bundle.
VERSION ?= 0.0.1

# List the GOOS and GOARCH to build
GO_LDFLAGS_STATIC="-s -w $(CTIMEVAR) -extldflags -static"

.DEFAULT_GOAL := help

##@ Build

build: ## Builds binaries for all supported platforms, linux, darwin
	go build -ldflags=${GO_LDFLAGS_STATIC} -o="$(BUILDDIR)/$(NAME)" main.go

##@ Testing

test: ## runs all tests
	go test -count=1 ./...

clean: ## Runs go clean
	go clean -i

lint: ## Runs golangci-lint
	golangci-lint run ./...

##@ Docker

docker_image: ## Creates a docker image for bundler. Requires `image` and `version` variables on command line
	docker build -t $(image):$(version) .

##@ Utilities

help:  ## Display this help. Thanks to https://www.thapaliya.com/en/writings/well-documented-makefiles/
ifeq ($(OS),Windows_NT)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-40s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
else
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-40s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
endif
