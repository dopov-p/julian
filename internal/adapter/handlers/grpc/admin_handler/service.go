package admin_handler

import (
	"context"

	admin "github.com/dopov-p/julian/internal/pb/admin/api"
	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
)

type (
	Service struct {
		admin.UnimplementedAdminServer

		cellUseCase cellUseCase
	}

	cellUseCase interface {
		Create(ctx context.Context, req *cell_usecase.CreateReq) error
	}
)

func NewService(cellUseCase cellUseCase) *Service {
	return &Service{
		cellUseCase: cellUseCase,
	}
}
