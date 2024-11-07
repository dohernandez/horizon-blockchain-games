GO ?= go

UNIT_TEST_COUNT ?= 2

# Override in Makefile to control unit test path.
UNIT_TEST_PATH ?= .

# Override in Makefile to control unit test package to exclude.
# To exclude multiple packages, use pipe (|) as separator.
EXCLUDE_TEST_PACKAGES ?=

# Override in Makefile to disable malformed LC_DYSYMTAB check during test.
DISABLE_MALFORMED_LC_DYSYMTAB ?=

## Run unit tests
test-unit:
	@echo "Running unit tests."
	@if [ ! -z "$(EXCLUDE_TEST_PACKAGES)" ]; then echo "Warning: Excluding package $(EXCLUDE_TEST_PACKAGES) from tests."; fi
	@CGO_ENABLED=1 $(GO) test -short -coverprofile=unit.coverprofile -covermode=atomic -race $(DISABLE_MALFORMED_LC_DYSYMTAB) `go list $(UNIT_TEST_PATH)/... | (! [ -z "$(EXCLUDE_TEST_PACKAGES)" ] && grep -vE "$(EXCLUDE_TEST_PACKAGES)" || cat)`


## Run unit tests multiple times, use `UNIT_TEST_COUNT=10 make test-unit-multi` to control count
test-unit-multi:
	@echo "Running unit tests ${UNIT_TEST_COUNT} times."
	@if [ ! -z "$(EXCLUDE_TEST_PACKAGES)" ]; then echo "Warning: Excluding package $(EXCLUDE_TEST_PACKAGES) from tests."; fi
	@CGO_ENABLED=1 $(GO) test -short -coverprofile=unit.coverprofile -count $(UNIT_TEST_COUNT) -covermode=atomic -race $(DISABLE_MALFORMED_LC_DYSYMTAB) `go list $(UNIT_TEST_PATH)/... | (! [ -z "$(EXCLUDE_TEST_PACKAGES)" ] && grep -vE "$(EXCLUDE_TEST_PACKAGES)" || cat)`

.PHONY: test-unit test-unit-multi
