package service_provider

import (
	"context"
	"strconv"

	"github.com/dopov-p/julian/internal/adapter/handlers/admin_handler"
	"github.com/dopov-p/julian/internal/adapter/handlers/cell_handler"
	admin "github.com/dopov-p/julian/internal/pb/admin/api"
	cell "github.com/dopov-p/julian/internal/pb/cell/api"
	"google.golang.org/grpc"
)

func (s *ServiceProvider) GetGRPCServer(ctx context.Context) interface {
	Start(ctx context.Context) error
	Stop()
} {
	if s.grpcServer == nil {
		cfg := s.GetConfig()
		port, err := strconv.Atoi(cfg.Server.Port)
		if err != nil {
			port = 8080
		}

		grpcServer := NewServer(port)
		grpcServerInstance := grpcServer.GetGRPCServer()

		// Register Admin service
		adminService := admin_handler.NewService(s.getCellUseCase(ctx))
		RegisterAdminServer(grpcServerInstance, adminService)

		// Register Cell service
		cellService := cell_handler.NewService(s.getCellUseCase(ctx))
		RegisterCellServer(grpcServerInstance, cellService)

		s.grpcServer = grpcServer
	}

	return s.grpcServer
}

// RegisterAdminServer registers Admin service with gRPC server
func RegisterAdminServer(s *grpc.Server, srv *admin_handler.Service) {
	admin.RegisterAdminServer(s, srv)
}

// RegisterCellServer registers Cell service with gRPC server
func RegisterCellServer(s *grpc.Server, srv *cell_handler.Service) {
	cell.RegisterCellServer(s, srv)
}
