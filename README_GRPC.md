# gRPC Setup Guide

## Prerequisites

1. **Install Protocol Buffers Compiler (protoc)**
   - macOS: `brew install protobuf`
   - Linux: `apt-get install protobuf-compiler` or download from [protobuf releases](https://github.com/protocolbuffers/protobuf/releases)
   - Windows: Download from [protobuf releases](https://github.com/protocolbuffers/protobuf/releases)
   - Verify: `protoc --version`

2. **Install Go plugins**:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

   Make sure `$GOPATH/bin` or `$HOME/go/bin` is in your `$PATH`.

## Generate gRPC Code

Run one of the following commands:

```bash
# Using Makefile (recommended)
make generate

# Or using script directly
./scripts/generate_proto.sh

# Or step by step
make install-tools
make proto
```

This will generate Go code in:
- `internal/pb/admin/` - Admin service code
- `internal/pb/cell/` - Cell service code

## After Generation

After running the generation, update the following files to use generated types:

1. **`internal/adapter/grpc/register.go`**:
   - Add imports: `"github.com/dopov-p/julian/internal/pb/admin"` and `"github.com/dopov-p/julian/internal/pb/cell"`
   - Replace function bodies with: `admin.RegisterAdminServer(s, srv)` and `cell.RegisterCellServer(s, srv)`

2. **`internal/adapter/grpc/admin_handler.go`**:
   - Uncomment `admin.UnimplementedAdminServer` embed
   - Update method signatures to use `*admin.CreateCellRequest` and `*admin.MarkCellDeletedRequest`
   - Implement the actual logic

3. **`internal/adapter/grpc/cell_handler.go`**:
   - Uncomment `cell.UnimplementedCellServer` embed
   - Update method signatures to use `*cell.DevastateCellContentsRequest` and `*cell.FullDevastateCellContentsRequest`
   - Uncomment and update the implementation code

## Running the Server

After updating the handlers, the gRPC server will start automatically when you run:

```bash
go run cmd/main.go
```

The server will listen on the port specified in `SERVER_PORT` environment variable (default: 8080).

## Testing with grpcurl

Install grpcurl:
```bash
# macOS
brew install grpcurl

# Or download from https://github.com/fullstorydev/grpcurl/releases
```

After starting the server, you can test it:

```bash
# List services
grpcurl -plaintext localhost:8080 list

# List methods for a service
grpcurl -plaintext localhost:8080 list cell.Cell

# Call a method
grpcurl -plaintext -d '{"data":{"name":"test","cell_contents":{"sku":"SKU123","quantity":10}}}' \
  localhost:8080 cell.Cell/DevastateCellContents
```
