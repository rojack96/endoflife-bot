#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

# Default version if none provided
DEFAULT_VERSION="v0.1.0"
ENV_FILE=".env"

# Priority: environment variable VERSION > .env file VERSION > DEFAULT_VERSION
VERSION="${VERSION:-}"

if [ -z "$VERSION" ] && [ -f "$ENV_FILE" ]; then
  # parse VERSION from .env (supports VERSION=... or VERSION="..." or VERSION='...')
  VERSION=$(grep -E '^\s*VERSION\s*=' "$ENV_FILE" | tail -n1 | sed -E 's/^\s*VERSION\s*=\s*//' | tr -d '"' | tr -d "'" | xargs)
fi

if [ -z "$VERSION" ]; then
  VERSION="$DEFAULT_VERSION"
fi

OUT="endoflife-${VERSION}"

echo "Building project..."
echo "Version: $VERSION"
echo "Output binary: $OUT"

# build (assumes main package is in the current directory)
go build -o "$OUT" .

echo "Build complete: $(pwd)/$OUT"