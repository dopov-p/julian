package cell_handler

import (
	"context"

	cell "github.com/dopov-p/julian/internal/pb/cell/api"
	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
)

type (
	Service struct {
		cell.UnimplementedCellServer

		cellRepo cellRepo

		cellUseCase cellUseCase
	}

	cellRepo interface {
		FullDevastate(ctx context.Context, name string) error
	}

	cellUseCase interface {
		DevastateContents(ctx context.Context, req *cell_usecase.DevastateContentsReq) error
	}
)

func NewService(cellRepo cellRepo, cellUseCase cellUseCase) *Service {
	return &Service{
		cellRepo:    cellRepo,
		cellUseCase: cellUseCase,
	}
}
