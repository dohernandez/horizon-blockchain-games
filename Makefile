GO ?= go

# Optional configuration to pinpoint the version of the tools.
# GCI_VERSION ?= "0.13.5"
# GOFUMPT_VERSION ?= "v0.7.0"
# GOLANGCI_LINT_VERSION ?= "v1.61.0"
# MOCKERY_VERSION ?= "2.46.3"

-include pkg/makefiles/main.mk
-include pkg/makefiles/lint.mk
-include pkg/makefiles/test-unit.mk
-include pkg/makefiles/mockery.mk
-include pkg/makefiles/build.mk

.PHONY: tools test generate check

## Install all require tools to work with the project
tools: install-mockery install-linter install-gci install-gofumpt

## Run tests
test: test-unit

generate:
	@echo "Generating ..."
	@go generate ./...

## Check test and linter the change.
check: lint test