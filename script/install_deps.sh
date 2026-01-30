#!/bin/sh
# Install protobuf compiler and Go gRPC/protobuf plugins if not already present.
# Idempotent: skips install when binaries exist. Run from repo root or anywhere.

set -e

# --- protoc ---
if command -v protoc >/dev/null 2>&1; then
  echo "protoc already installed: $(protoc --version 2>/dev/null || true)"
else
  echo "Installing protoc..."
  if command -v apt-get >/dev/null 2>&1; then
    sudo apt-get update -qq && sudo apt-get install -y protobuf-compiler
  elif command -v apk >/dev/null 2>&1; then
    sudo apk add --no-cache protobuf-dev
  elif command -v brew >/dev/null 2>&1; then
    brew install protobuf
  else
    echo "No supported package manager (apt-get/apk/brew). Install protoc manually." >&2
    exit 1
  fi
fi

# Ensure Go tools are on PATH (common for go install)
GOPATH="${GOPATH:-$HOME/go}"
export PATH="$PATH:$GOPATH/bin"

# --- protoc-gen-go ---
if command -v protoc-gen-go >/dev/null 2>&1; then
  echo "protoc-gen-go already installed"
else
  echo "Installing protoc-gen-go..."
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# --- protoc-gen-go-grpc ---
if command -v protoc-gen-go-grpc >/dev/null 2>&1; then
  echo "protoc-gen-go-grpc already installed"
else
  echo "Installing protoc-gen-go-grpc..."
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

echo "Done. protoc and Go plugins are ready."
