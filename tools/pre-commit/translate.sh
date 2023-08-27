#!/usr/bin/env bash

set -e -o pipefail

if ! command -v go version &> /dev/null ; then
    echo "go not installed or available in the PATH" >&2
    exit 1
fi

export PATH="$(go env GOPATH)/bin:$PATH"
export PATH="$(go env GOROOT)/bin:$PATH"

# if file bin/optimizer not exists, build it
if ! command -v content-translator &> /dev/null ; then
  go install github.com/chloyka/chloyka.com/tools/pre-commit/content-translator@latest
fi

content-translator "$@"
