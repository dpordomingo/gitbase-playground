# Package configuration
PROJECT := gitbase-playground
COMMANDS := cmd/server
DEPENDENCIES := \
	github.com/golang/dep/cmd/dep \
	github.com/jteeuwen/go-bindata \
	github.com/golang/lint/golint
GO_LINTABLE_PACKAGES := $(shell go list ./... | grep -v '/vendor/')

# Tools
GODEP := dep
GOLINT := golint
GOVET := go vet
BINDATA := go-bindata

all:

# Including ci Makefile
CI_REPOSITORY ?= https://github.com/src-d/ci.git
CI_BRANCH ?= v1
CI_PATH ?= $(shell pwd)/.ci
MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	@git clone --quiet --depth 1 -b $(CI_BRANCH) $(CI_REPOSITORY) $(CI_PATH);
-include $(MAKEFILE)

# Makefile.main::dependencies -> Makefile.main::$(DEPENDENCIES)
dependencies: front-dependencies back-dependencies
	exit 0;

# Makefile.main::build -> Makefile.main::$(COMMANDS)
build: front-build back-build

# Makefile.main::test
test: back-test front-test

converage: test-coverage codecov # from Makefile.main

lint: back-lint front-lint


# Backend

back-dependencies:
	$(GODEP) ensure

back-build: back-bindata

back-bindata:
	$(BINDATA) \
		-pkg assets \
		-o ./server/assets/asset.go \
		build/public/*

back-lint: $(GO_LINTABLE_PACKAGES)
$(GO_LINTABLE_PACKAGES):
	$(GOLINT) $@
	$(GOVET) $@

back-start:
	go run cli/server/server.go


# Frontend
yarn_production ?= true

front-dependencies-development:
	$(MAKE) front-dependencies yarn_production=false

front-dependencies:
	$(YARN) install --production=$(yarn_production)

front-test: front-dependencies-development
	$(YARN) test

front-lint: front-dependencies-development
	$(YARN) lint

front-build: front-dependencies
	$(YARN) build

front-start: front-dependencies
	$(YARN) start


# ALL

prepare-build: | front-build back-build bindata

validate-commit: | back-dependencies no-changes-in-commit

build-app: | prepare-build packages

## Compiles the assets, and serve the tool through its API
serve: | front-build back-build gorun
