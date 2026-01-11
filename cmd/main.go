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

	sp := service_provider.NewServiceProvider()

	// Initialize database cluster
	dbCluster := sp.GetDbCluster(ctx)

	// Initialize and start gRPC server
	grpcServer := sp.GetGRPCServer(ctx)
	err := grpcServer.Start(ctx)
	if err != nil {
		cancel()          // Cancel context before exit
		dbCluster.Close() // Ensure database is closed before exit
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	// Close database when server returns (server.Start blocks until ctx.Done)
	// Context is already cancelled at this point (server.Start blocks until ctx.Done)
	cancel()
	dbCluster.Close()
}
