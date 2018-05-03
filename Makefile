# Package configuration
PROJECT := gitbase-playground
COMMANDS := cmd/server
DEPENDENCIES := \
	github.com/golang/dep/cmd/dep \
	github.com/jteeuwen/go-bindata \
	github.com/golang/lint/golint
GO_LINTABLE_PACKAGES := $(shell go list ./... | grep -v '/vendor/')

assets := ./server/assets/asset.go
assets_back := $(assets).bak

# Tools
GODEP := dep
GOLINT := golint
GOVET := go vet
BINDATA := go-bindata
DIFF := diff

all:

# Including ci Makefile
CI_REPOSITORY ?= https://github.com/src-d/ci.git
CI_BRANCH ?= v1
CI_PATH ?= $(shell pwd)/.ci
MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	@git clone --quiet --depth 1 -b $(CI_BRANCH) $(CI_REPOSITORY) $(CI_PATH);
-include $(MAKEFILE)

# Makefile.main::dependencies -> Makefile.main::$(DEPENDENCIES) -> this::dependencies
dependencies: | front-dependencies back-dependencies exit

# Makefile.main::test -> this::test
test: front-test

# this::build -> Makefile.main::build -> Makefile.main::$(COMMANDS)
build: prepare-build
	@echo

prepare-build: | front-build back-build

coverage: | test-coverage codecov

lint: | back-lint front-lint


# Backend

back-dependencies:
	$(GODEP) ensure

back-build: back-bindata

back-bindata:
	$(BINDATA) \
		-pkg assets \
		-o $(assets) \
		build/public/*

back-lint: $(GO_LINTABLE_PACKAGES)
$(GO_LINTABLE_PACKAGES):
	$(GOLINT) $@
	$(GOVET) $@

back-start:
	go run cmd/server/main.go

back-ensure-assets-proxy:
	$(DIFF) $(assets) $(assets_back) || exit 1


# Frontend
yarn_production ?= true

front-dependencies-development:
	echo 'SKIP. no frontend dependencies to install'
	#$(MAKE) front-dependencies yarn_production=false

front-dependencies:
	echo 'SKIP. no frontend dependencies to install'
	#$(YARN) install --production=$(yarn_production)

front-test: front-dependencies-development
	echo 'SKIP. no frontend tests to run'
	#$(YARN) test

front-lint: front-dependencies-development
	echo 'SKIP. no frontend linters to run'
	#$(YARN) lint

front-build:
	mkdir -p build/public
	cp public/index.html build/public/index.html
	cp public/secondary.html build/public/secondary.html
	#$(YARN) build

front-start: front-dependencies
	echo 'SKIP. no runnable frontend'
	#$(YARN) start

front-fix-lint-errors:
	echo 'SKIP. no fixable code'
	#$(YARN) fix-lint-errors

# ALL

exit:
	exit 0;

validate-commit: | back-dependencies back-ensure-assets-proxy front-fix-lint-errors no-changes-in-commit

build-app: | prepare-build packages

## Compiles the assets, and serve the tool through its API
serve: | front-build back-build gorun

