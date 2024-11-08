GO ?= go

# Override in app main.mk to add custom ldflags, example BUILD_LDFLAGS="-s -w".
BUILD_LDFLAGS ?= ""

# Override in app main.mk to control build target, example BUILD_PKG=./cmd/my-app.
BUILD_PKG ?= ./...

# Override in app main.mk to control build artifact destination.
BUILD_DIR ?= ./bin

## Build Linux binary
build-linux:
	@echo "Building Linux AMD64 binary, GOFLAGS: $(GOFLAGS)"
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(shell bash $(MAKEFILES_PATH)/version-ldflags.sh && echo $(BUILD_LDFLAGS))" -o $(BUILD_DIR)/ $(BUILD_PKG)

## Build macOS intel binary
build-darwin-amd:
	@echo "Building macOS INTEL (Darwin AMD64) binary, GOFLAGS: $(GOFLAGS)"
	@GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(shell bash $(MAKEFILES_PATH)/version-ldflags.sh && echo $(BUILD_LDFLAGS))" -o $(BUILD_DIR) $(BUILD_PKG)

## Build macOS Apple M1 binary
build-darwin-arm:
	@echo "Building macOS Apple M1 (Darwin ARM64) binary, GOFLAGS: $(GOFLAGS)"
	@GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(shell bash $(MAKEFILES_PATH)/version-ldflags.sh && echo $(BUILD_LDFLAGS))" -o $(BUILD_DIR) $(BUILD_PKG)

## Build binary
build:
	@echo "Building binary, GOFLAGS: $(GOFLAGS)"
	@$(GO) build -ldflags "$(shell bash $(MAKEFILES_PATH)/version-ldflags.sh && echo $(BUILD_LDFLAGS))" -o $(BUILD_DIR)/ $(BUILD_PKG)


.PHONY: build-linux build-darwin-amd build-darwin-arm build