.PHONY: proto generate install-tools check-protoc install-linters check-golangci-lint install-pre-commit check-pre-commit lint lint-fix format

# Check if protoc is installed
check-protoc:
	@which protoc > /dev/null || (echo "Error: protoc is not installed. Please install Protocol Buffers compiler:" && \
		echo "  macOS: brew install protobuf" && \
		echo "  Linux: apt-get install protobuf-compiler" && \
		echo "  Or download from: https://grpc.io/docs/protoc-installation/" && exit 1)

# Install required Go tools for proto generation and code formatting
install-tools:
	@echo "Installing protoc-gen-go and protoc-gen-go-grpc..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Installing goimports for code formatting..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Go tools installed successfully!"
	@echo "Note: Make sure protoc is installed (see check-protoc target)"

# Generate gRPC code from proto files
proto: check-protoc
	@echo "Generating gRPC code from proto files..."
	@mkdir -p internal/pb/admin internal/pb/cell
	@echo "Generating code for admin.proto..."
	@protoc \
		--go_out=internal/pb/admin \
		--go_opt=paths=source_relative \
		--go-grpc_out=internal/pb/admin \
		--go-grpc_opt=paths=source_relative \
		api/admin.proto
	@echo "Generating code for cell.proto..."
	@protoc \
		--go_out=internal/pb/cell \
		--go_opt=paths=source_relative \
		--go-grpc_out=internal/pb/cell \
		--go-grpc_opt=paths=source_relative \
		api/cell.proto
	@echo "gRPC code generated successfully!"

# Generate and install tools
generate: install-tools proto

# Check if golangci-lint is installed
check-golangci-lint:
	@which golangci-lint > /dev/null || (echo "Error: golangci-lint is not installed. Run 'make install-linters' to install it." && exit 1)

# Install golangci-lint
install-linters:
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest
	@echo "golangci-lint installed successfully!"
	@golangci-lint --version

# Check if pre-commit is installed
check-pre-commit:
	@which pre-commit > /dev/null || (echo "Error: pre-commit is not installed. Install it with:" && \
		echo "  macOS: brew install pre-commit" && \
		echo "  Linux: pip install pre-commit" && \
		echo "  Or: pip3 install pre-commit" && exit 1)

# Install pre-commit hooks
install-pre-commit: check-pre-commit
	@echo "Installing pre-commit hooks..."
	@pre-commit install
	@pre-commit install --hook-type pre-push
	@echo "Pre-commit hooks installed successfully!"

# Run linters
lint: check-golangci-lint
	@echo "Running golangci-lint..."
	@golangci-lint run

# Run linters with auto-fix
lint-fix: check-golangci-lint
	@echo "Running golangci-lint with auto-fix..."
	@golangci-lint run --fix

# Run pre-commit on all files
format: check-pre-commit
	@echo "Running pre-commit on all files..."
	@PATH="$$(go env GOPATH)/bin:$$PATH" pre-commit run --all-files
