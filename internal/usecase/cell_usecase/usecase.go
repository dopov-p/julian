package cell_usecase

import (
	"context"

	"github.com/dopov-p/julian/internal/domain/dto"
	"github.com/dopov-p/julian/internal/domain/model"
)

type (
	UseCase struct {
		cellRepo  cellRepo
		generator generator
		txManager txManager
	}

	cellRepo interface {
		Create(ctx context.Context, cell *model.Cell) error
		GetContentsByName(ctx context.Context, name string) (*dto.GetContentsByNameResponse, error)
		UpdateContents(ctx context.Context, req dto.UpdateContentsRequest) error
		GetByName(ctx context.Context, name string, isActive bool) (*model.Cell, error)
		MarkDeleted(ctx context.Context, name string) error
		MarkActivated(ctx context.Context, name string) error
	}

	generator interface {
		GenerateUUID() (string, error)
	}
	txManager interface {
		WithTx(ctx context.Context, fn func(context.Context) error) error
	}
)

func NewUseCase(cellRepo cellRepo, generator generator, txManager txManager) *UseCase {
	return &UseCase{
		cellRepo:  cellRepo,
		generator: generator,
		txManager: txManager,
	}
}
