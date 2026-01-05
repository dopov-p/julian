package service_provider

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	port       int
}

func NewServer(port int) *Server {
	// Create gRPC server with default options
	grpcServer := grpc.NewServer()

	// Enable reflection for development (can be disabled in production)
	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		port:       port,
	}
}

// GetGRPCServer returns the underlying grpc.Server for direct registration
func (s *Server) GetGRPCServer() *grpc.Server {
	return s.grpcServer
}

// Start starts the gRPC server and blocks until context is cancelled
func (s *Server) Start(ctx context.Context) error {
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

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	log.Println("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
}
