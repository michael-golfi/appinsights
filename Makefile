PROJECT_NAME := "appinsights"
IMAGE_NAME := "michael.golfi/appinsights"
PKG := "gitlab.com/michael.golfi/appinsights"
TAG ?= "latest"

PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all install build lint test race msan coverage coverhtml deploy clean help

all: build

#
# Build
#
install: ## Get the dependencies
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure
	@go get -u github.com/golang/lint/golint

build: #dep ## Build the binary file
	GOOS=linux go build -i -v $(PKG)

#
# Test
# 
lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	GOOS=linux go test -short ${PKG_LIST}

race: #dep ## Run data race detector
	CGO_ENABLED=1 GOOS=linux go test -race -short ${PKG_LIST}

msan: #dep ## Run memory sanitizer
	CGO_ENABLED=1 GOOS=linux go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	@chmod +x tools/coverage.sh
	GOOS=linux ./tools/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	GOOS=linux ./tools/coverage.sh html;

#
# Deploy
#
deploy:
	./scripts/build.sh
	#@docker plugin push $(IMAGE_NAME):$(TAG)

#
# Help and Teardown
#
clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'