package service_provider

import (
	"context"

	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
)

func (s *ServiceProvider) getCellUseCase(ctx context.Context) *cell_usecase.UseCase {
	if s.cellUseCase == nil {
		s.cellUseCase = cell_usecase.NewUseCase(
			s.getCellRepo(ctx),
			s.getGenerator(),
			s.getTxManager(ctx),
		)
	}

	return s.cellUseCase
}
