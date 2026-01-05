package cell_repo

import (
	"context"
	"time"

	"github.com/dopov-p/julian/internal/config"
	"github.com/dopov-p/julian/internal/pkg"
)

const (
	tableName        = "cells"
	defaultBatchSize = uint64(50)
)

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

func (r *Repo) getConn(ctx context.Context) pkg.Querier {
	if tx := pkg.GetTx(ctx); tx != nil {
		return tx
	}

	return r.cluster.Conn
}
