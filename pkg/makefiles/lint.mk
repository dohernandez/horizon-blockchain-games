GO ?= go

# Override in Makefile to control lint path.
LINT_PATH ?= .

# Override in Makefile to control cmd lint path.
CMD_LINT_PATH ?= ./*/cmd

## Check/install golangci-lint tool
install-linter:
	@GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION) bash $(MAKEFILES_PATH)/install-linter.sh

## Check/install gci tool
install-gci:
	@GCI_VERSION=$(GCI_VERSION) bash $(MAKEFILES_PATH)/install-gci.sh

## Check/install gofumpt tool
install-gofumpt:
	@GOFUMPT_VERSION=$(GOFUMPT_VERSION) bash $(MAKEFILES_PATH)/install-gofumpt.sh

## Check with golangci-lint
lint:
	@LINT_PATH=$(LINT_PATH) CMD_LINT_PATH=$(CMD_LINT_PATH) bash $(MAKEFILES_PATH)/lint.sh

## Apply goimports and gofmt
fix-lint:
	@bash $(MAKEFILES_PATH)/fix.sh

.PHONY: install-linter install-gci install-gofumpt lint fix-lint