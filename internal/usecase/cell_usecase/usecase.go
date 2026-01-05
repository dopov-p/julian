package cell_usecase

import (
	"context"

	"github.com/dopov-p/julian/internal/config"
	"github.com/dopov-p/julian/internal/domain/dto"
	"github.com/dopov-p/julian/internal/pkg"
)

type (
	UseCase struct {
		cellRepo cellRepo
		cluster  *config.Cluster
	}

	cellRepo interface {
		GetContentsByName(ctx context.Context, name string) (*dto.GetContentsByNameResponse, error)
		UpdateContents(ctx context.Context, req dto.UpdateContentsRequest) error
	}
)

func NewUseCase(cellRepo cellRepo, cluster *config.Cluster) *UseCase {
	return &UseCase{
		cellRepo: cellRepo,
		cluster:  cluster,
	}
}

func (u *UseCase) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return pkg.WithTx(ctx, u.cluster.Conn, fn)
}
