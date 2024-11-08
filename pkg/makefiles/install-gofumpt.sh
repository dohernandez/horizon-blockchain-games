#!/usr/bin/env bash

set -e -o pipefail

[ -z "$GO" ] && GO=go

# Override in Makefile to control gofumpt version.
[ -z "$GOFUMPT_VERSION" ] && GOFUMPT_VERSION="v0.7.0"

# detecting GOPATH and removing trailing "/" if any
GOPATH="$(go env GOPATH)"
GOPATH=${GOPATH%/}

# adding GOBIN to PATH
[[ ":$PATH:" != *"$GOPATH/bin"* ]] && PATH=$PATH:"$GOPATH"/bin

# checking if gofumpt is available and it is the version specify
# gofumpt is a drop-in replacement for gofmt with stricter formatting: https://github.com/mvdan/gofumpt
if ! command -v gofumpt > /dev/null; then \
  echo ">> Installing gofumpt $GOFUMPT_VERSION..."; \
  $GO install mvdan.cc/gofumpt@"$GOFUMPT_VERSION";
else
  VERSION_INSTALLED="$(gofumpt --version | awk '{print $1}')"

  if [ "${VERSION_INSTALLED}" != "${GOFUMPT_VERSION}" ]; then \
    echo ">> Updating gofumpt from ${VERSION_INSTALLED} to $GOFUMPT_VERSION..."; \
    $GO install mvdan.cc/gofumpt@"$GOFUMPT_VERSION";
  fi
fi