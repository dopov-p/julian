#!/bin/bash

# Script to generate gRPC code from proto files
# Usage: ./scripts/generate_proto.sh

set -e

echo "Generating gRPC code from proto files..."

# Create output directories
mkdir -p internal/pb/admin
mkdir -p internal/pb/cell

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Please install Protocol Buffers compiler:"
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt-get install protobuf-compiler"
    echo "  Or download from: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Check if protoc-gen-go is in PATH
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Error: protoc-gen-go is not in PATH"
    echo "Please run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# Check if protoc-gen-go-grpc is in PATH
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Error: protoc-gen-go-grpc is not in PATH"
    echo "Please run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# Generate code for admin.proto
echo "Generating code for admin.proto..."
protoc \
    --go_out=internal/pb/admin \
    --go_opt=paths=source_relative \
    --go-grpc_out=internal/pb/admin \
    --go-grpc_opt=paths=source_relative \
    api/admin.proto

# Generate code for cell.proto
echo "Generating code for cell.proto..."
protoc \
    --go_out=internal/pb/cell \
    --go_opt=paths=source_relative \
    --go-grpc_out=internal/pb/cell \
    --go-grpc_opt=paths=source_relative \
    api/cell.proto

echo "gRPC code generated successfully!"
