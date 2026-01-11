#!/bin/bash
# Wrapper script for golangci-lint to ensure GOPATH/bin is in PATH
export PATH="$(go env GOPATH)/bin:$PATH"
exec golangci-lint run --fix --timeout=5m "$@"
