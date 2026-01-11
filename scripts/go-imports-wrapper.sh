#!/bin/bash
# Wrapper script for go-imports to ensure GOPATH/bin is in PATH
export PATH="$(go env GOPATH)/bin:$PATH"
exec goimports -w "$@"
