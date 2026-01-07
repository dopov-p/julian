package service_provider

import (
	"context"
	"log"

	"github.com/dopov-p/julian/internal/config"
)

func (s *ServiceProvider) GetDbCluster(ctx context.Context) *config.Cluster {
	if s.dbCluster == nil {
		dbCluster, err := config.NewCluster(ctx)
		if err != nil {
			log.Fatal(err)
		}

		s.dbCluster = dbCluster
	}

	return s.dbCluster
}
