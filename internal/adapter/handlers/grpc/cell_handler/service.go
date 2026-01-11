package cell_handler

import (
	"context"

	cell "github.com/dopov-p/julian/internal/pb/cell/api"
	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
)

type (
	Service struct {
		cell.UnimplementedCellServer
		cellUseCase cellUseCase
	}

	cellUseCase interface {
		DevastateContents(ctx context.Context, req cell_usecase.DevastateContentsReq) error
	}
)

func NewService(cellUseCase cellUseCase) *Service {
	return &Service{
		cellUseCase: cellUseCase,
	}
}
