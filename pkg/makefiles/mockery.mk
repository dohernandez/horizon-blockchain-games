GO ?= go

# Override in Makefile to control mockery version.
MOCKERY_VERSION ?= "2.46.3"

## Check/install mockery tool
install-mockery:
	@MOCKERY_VERSION=$(MOCKERY_VERSION) bash $(MAKEFILES_PATH)/install-mockery.sh

.PHONY: install-mockery