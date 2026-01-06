.PHONY: proto generate install-tools check-protoc

# Check if protoc is installed
check-protoc:
	@which protoc > /dev/null || (echo "Error: protoc is not installed. Please install Protocol Buffers compiler:" && \
		echo "  macOS: brew install protobuf" && \
		echo "  Linux: apt-get install protobuf-compiler" && \
		echo "  Or download from: https://grpc.io/docs/protoc-installation/" && exit 1)

# Install required Go tools for proto generation
install-tools:
	@echo "Installing protoc-gen-go and protoc-gen-go-grpc..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
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
