#!/usr/bin/env bash

branch=$GITHUB_HEAD_REF
[[ -z "$branch" ]] && branch=${GITHUB_REF#refs/heads/}

# skip branch/revision with missing git repo.
if [[ -d .git ]] || git rev-parse --git-dir > /dev/null 2>&1; then
    [[ -z "$branch" ]] && branch=$(git symbolic-ref HEAD 2>/dev/null)
    [[ -z "$VERSION" ]] && VERSION=$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
    [[ -z "$VERSION" ]] && VERSION=branch
    revision=$(git log -1 --pretty=format:"%H" 2>/dev/null)
fi

build_user="$USER"
build_date=$(date +%FT%T%Z)

# gh_repo is the variable available in gh-repo.sh that contains the long repo name. Example: github.com/dohernandez/horizon-blockchain-games
BASEDIR=$(dirname "$0")
source "$BASEDIR"/gh-repo.sh

version_pkg="github.com/dohernandez/horizon-blockchain-games/version"

echo -X "$version_pkg".version="$VERSION" -X "$version_pkg".branch="$branch" -X "$version_pkg".revision="$revision" -X "$version_pkg".buildUser="$build_user" -X "$version_pkg".buildDate="$build_date"