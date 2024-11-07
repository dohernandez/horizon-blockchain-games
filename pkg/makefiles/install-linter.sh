#!/usr/bin/env bash

[ -z "$GO" ] && GO=go

# Override in Makefile to control gofumpt version.
[ -z "$GOLANGCI_LINT_VERSION" ] && GOLANGCI_LINT_VERSION="1.61.0"

[[ $GOLANGCI_LINT_VERSION == v* ]] && GOLANGCI_LINT_VERSION="${GOLANGCI_LINT_VERSION:1}"

install_golangci_lint () {
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /tmp/golangci-lint-"$1" v"$1" \
            && mv /tmp/golangci-lint-"$1"/golangci-lint "$GOPATH"/bin/golangci-lint
}

# detecting GOPATH and removing trailing "/" if any
GOPATH="$(go env GOPATH)"
GOPATH=${GOPATH%/}

# adding GOBIN to PATH
[[ ":$PATH:" != *"$GOPATH/bin"* ]] && PATH=$PATH:"$GOPATH"/bin

# checking if golangci-lint is available and it is the version specify
if ! command -v golangci-lint >/dev/null; then
  echo ">> Installing golangci-lint v$GOLANGCI_LINT_VERSION...";
  install_golangci_lint "$GOLANGCI_LINT_VERSION"
else
    VERSION_INSTALLED="$(golangci-lint --version | sed 's/[^0-9.]*\([0-9.]*\).*/\1/')"
    if [ "${VERSION_INSTALLED}" != "${GOLANGCI_LINT_VERSION}" ]; then \
      echo ">> Updating golangci-lint from v"${VERSION_INSTALLED}" to v$GOLANGCI_LINT_VERSION..."; \
      install_golangci_lint "$GOLANGCI_LINT_VERSION"
    fi
fi