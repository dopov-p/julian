#!/bin/bash
# Wrapper script for golines to ensure GOPATH/bin is in PATH
export PATH="$(go env GOPATH)/bin:$PATH"
exec golines -w "$@"
