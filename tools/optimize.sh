#!/usr/bin/env bash

set -e -o pipefail

if ! command -v go version &> /dev/null ; then
    echo "go not installed or available in the PATH" >&2
    exit 1
fi
# if file bin/optimizer not exists, build it
if ![ -f bin/optimizer ]; then
    cd image-optimizer && go build -o ../bin/optimizer && cd ..
fi

./bin/optimizer -dir=$1
