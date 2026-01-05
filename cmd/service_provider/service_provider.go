package service_provider

import (
	"context"

	"github.com/dopov-p/julian/internal/adapter/storage/postgres/cell_repo"
	"github.com/dopov-p/julian/internal/config"
	"github.com/dopov-p/julian/internal/pkg"
	"github.com/dopov-p/julian/internal/usecase/cell_usecase"
)

type ServiceProvider struct {
	dbCluster *config.Cluster
	config    *config.Config

	cellRepo *cell_repo.Repo

	cellUseCase *cell_usecase.UseCase

	timer     *pkg.Timer
	generator *pkg.Generator

	grpcServer interface {
		Start(ctx context.Context) error
		Stop()
	}
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) getTimer() *pkg.Timer {
	if s.timer == nil {
		s.timer = pkg.NewTimer()
	}

	return s.timer
}

func (s *ServiceProvider) getGenerator() *pkg.Generator {
	if s.generator == nil {
		s.generator = pkg.NewGenerator()
	}

	return s.generator
}

func (s *ServiceProvider) GetConfig() *config.Config {
	if s.config == nil {
		cfg, err := config.Load()
		if err != nil {
			panic(err)
		}
		s.config = cfg
	}

	return s.config
}
