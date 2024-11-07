#!/usr/bin/env bash

[ -z "$GO" ] && GO=go
[ -z "$LINT_PATH" ] && LINT_PATH="."

# detecting GOPATH and removing trailing "/" if any
GOPATH="$(go env GOPATH)"
GOPATH=${GOPATH%/}

# adding GOBIN to PATH
[[ ":$PATH:" != *"$GOPATH/bin"* ]] && PATH=$PATH:"$GOPATH"/bin

if ! command -v golangci-lint >/dev/null; then
  echo "golangci-lint is not installed"
fi

this_path=$(dirname "$0")

golangci_yml="./.golangci.yml"
if [ ! -f "./.golangci.yml" ]; then
  golangci_yml="$this_path"/.golangci.yml
fi

echo "Checking packages."
golangci-lint run -c "$golangci_yml" "$LINT_PATH"/... || exit 1