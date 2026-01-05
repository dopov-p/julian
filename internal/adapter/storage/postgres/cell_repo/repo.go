package cell_repo

import (
	"time"

	"github.com/dopov-p/julian/internal/config"
)

const tableName = "cells"

type (
	Repo struct {
		cluster *config.Cluster

		timer timer
	}

	timer interface {
		NowUTC() time.Time
	}
)

func NewRepo(cluster *config.Cluster, timer timer) *Repo {
	return &Repo{
		cluster: cluster,
		timer:   timer,
	}
}
