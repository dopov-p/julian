package cell_handler

import (
	cell "github.com/dopov-p/julian/internal/pb/cell/api"
	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
)

// Service implements the Cell gRPC service
type Service struct {
	cell.UnimplementedCellServer
	cellUseCase *cell_usecase.UseCase
}

func NewService(cellUseCase *cell_usecase.UseCase) *Service {
	return &Service{
		cellUseCase: cellUseCase,
	}
}
