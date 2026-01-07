package admin_handler

import (
	admin "github.com/dopov-p/julian/internal/pb/admin/api"
	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
)

// Service implements the Admin gRPC service
type Service struct {
	admin.UnimplementedAdminServer
	cellUseCase *cell_usecase.UseCase
}

func NewService(cellUseCase *cell_usecase.UseCase) *Service {
	return &Service{
		cellUseCase: cellUseCase,
	}
}
