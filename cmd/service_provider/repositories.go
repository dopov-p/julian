package service_provider

import (
	"context"

	"github.com/dopov-p/julian/internal/adapter/storage/postgres/cell_repo"
)

func (s *ServiceProvider) getCellRepo(ctx context.Context) *cell_repo.Repo {
	if s.cellRepo == nil {
		s.cellRepo = cell_repo.NewRepo(
			s.GetDbCluster(ctx),
			s.getTimer(),
		)
	}

	return s.cellRepo
}
