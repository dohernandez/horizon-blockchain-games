#!/usr/bin/env bash

[ -z "$GO" ] && GO=go

# Override in Makefile to control gci version.
[ -z "$GCI_VERSION" ] && GCI_VERSION="0.13.5"

# detecting GOPATH and removing trailing "/" if any
GOPATH="$(go env GOPATH)"
GOPATH=${GOPATH%/}

# adding GOBIN to PATH
[[ ":$PATH:" != *"$GOPATH/bin"* ]] && PATH=$PATH:"$GOPATH"/bin

# checking if gci is available and it is the version specify
# gci a tool that controls golang package import order and makes it always deterministic. https://github.com/daixiang0/gci
if ! command -v gci > /dev/null; then \
    echo ">> Installing gci v$GCI_VERSION..."; \
    $GO install -mod=mod github.com/daixiang0/gci@"v$GCI_VERSION";
else
  VERSION_INSTALLED="$(gci --version | cut -d' ' -f3)"
  if [ "${VERSION_INSTALLED}" != "${GCI_VERSION}" ]; then \
    echo ">> Updating gci form v"${VERSION_INSTALLED}" to v$GCI_VERSION..."; \
    $GO install -mod=mod github.com/daixiang0/gci@"v$GCI_VERSION";
  fi
fi