#!/usr/bin/env bash

set -e -o pipefail

if ! command -v go version &> /dev/null ; then
    echo "go not installed or available in the PATH" >&2
    exit 1
fi
# if file bin/optimizer not exists, build it
if ! command -v image-optimizer &> /dev/null ; then
  go install github.com/chloyka/chloyka.com/tools/image-optimizer@latest
fi

exec image-optimizer "$@"
