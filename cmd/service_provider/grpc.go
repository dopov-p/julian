package service_provider

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/dopov-p/julian/internal/adapter/handlers/grpc/admin_handler"
	"github.com/dopov-p/julian/internal/adapter/handlers/grpc/cell_handler"
	adminPb "github.com/dopov-p/julian/internal/pb/admin/api"
	cellPb "github.com/dopov-p/julian/internal/pb/cell/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s *ServiceProvider) GetGRPCServer(ctx context.Context) interface {
	Start(ctx context.Context) error
	Stop()
} {
	if s.grpcServer == nil {
		cfg := s.GetConfig()
		port, err := strconv.Atoi(cfg.GrpcServer.Port)
		if err != nil {
			port = 50051
		}

		grpcServer := NewGrpcServer(port)
		grpcServerInstance := grpcServer.GetGRPCServer()

		// Register Admin service
		adminService := admin_handler.NewService(
			s.getCellUseCase(ctx),
		)
		RegisterAdminServer(grpcServerInstance, adminService)

		// Register Cell service
		cellService := cell_handler.NewService(
			s.getCellRepo(ctx),
			s.getCellUseCase(ctx),
		)
		RegisterCellServer(grpcServerInstance, cellService)

		s.grpcServer = grpcServer
	}

	return s.grpcServer
}

// RegisterAdminServer registers Admin service with gRPC server.
func RegisterAdminServer(s *grpc.Server, srv *admin_handler.Service) {
	adminPb.RegisterAdminServer(s, srv)
}

// RegisterCellServer registers Cell service with gRPC server.
func RegisterCellServer(s *grpc.Server, srv *cell_handler.Service) {
	cellPb.RegisterCellServer(s, srv)
}

type GrpcServer struct {
	grpcServer *grpc.Server
	port       int
}

func NewGrpcServer(port int) *GrpcServer {
	// Create gRPC server with default options
	grpcServer := grpc.NewServer()

	// Enable reflection for development (can be disabled in production)
	reflection.Register(grpcServer)

	return &GrpcServer{
		grpcServer: grpcServer,
		port:       port,
	}
}

// GetGRPCServer returns the underlying grpc.Server for direct registration.
func (s *GrpcServer) GetGRPCServer() *grpc.Server {
	return s.grpcServer
}

// Start starts the gRPC server and blocks until context is cancelled.
func (s *GrpcServer) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.port, err)
	}

	log.Printf("gRPC server listening on port %d", s.port)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err = s.grpcServer.Serve(lis); err != nil {
			errChan <- fmt.Errorf("gRPC server failed: %w", err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		log.Println("Shutting down gRPC server...")
		s.grpcServer.GracefulStop()
		return nil
	case err = <-errChan:
		return err
	}
}

// Stop gracefully stops the gRPC server.
func (s *GrpcServer) Stop() {
	log.Println("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
}
