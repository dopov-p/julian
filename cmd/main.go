package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/dopov-p/julian/cmd/service_provider"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	sp := service_provider.NewServiceProvider()

	// Initialize database cluster
	dbCluster := sp.GetDbCluster(ctx)
	defer dbCluster.Close()

	// Initialize and start gRPC server
	grpcServer := sp.GetGRPCServer(ctx)
	if err := grpcServer.Start(ctx); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
