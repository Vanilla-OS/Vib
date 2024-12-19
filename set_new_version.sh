#!/usr/bin/env bash

set -eu

if [ $# -lt 1 ]; then
    echo "usag: $0 <version>"
    exit 1
fi

IFS='.' read -ra VERS <<< ${1##v}
if [ ${#VERS[@]} -lt 3 ]; then
    echo "invalid version number"
    exit 2
fi

sed -i "s|var Version = \"0.0.0\"|var Version = \"${1##v}\"|g" cmd/root.go

sed -i "s|0, 0, 0|${VERS[0]}, ${VERS[1]}, ${VERS[2]}|" core/loader.go

