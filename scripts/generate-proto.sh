#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

export PATH="${PATH}:${HOME}/go/bin"

if ! command -v buf >/dev/null 2>&1; then
  GOWORK=off go install github.com/bufbuild/buf/cmd/buf@v1.47.2
fi

if ! command -v protoc-gen-go >/dev/null 2>&1; then
  GOWORK=off go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
fi

buf generate
mkdir -p descriptorsets
buf build proto -o descriptorsets/control_plane.binpb
