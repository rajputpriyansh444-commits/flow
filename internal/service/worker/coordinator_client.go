package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/logger"
	"github.com/rajputpriyansh444-commits/flow/internal/runtime/workspacebundle"
	"github.com/rajputpriyansh444-commits/flow/internal/service/coordinator"
)

// NewCoordinatorClient creates a client using the worker's configured coordinators.
func NewCoordinatorClient(ctx context.Context, cfg *config.Config) (coordinator.Client, error) {
	if len(cfg.Worker.Coordinators) == 0 {
		return nil, fmt.Errorf("worker.coordinators is required")
	}

	clientConfig := coordinator.ConfigFromPeer(cfg.Core.Peer)
	clientConfig.WorkspaceBundleDir = workspacebundle.StoreDir(cfg.Paths.DataDir)
	if err := clientConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid coordinator client configuration: %w", err)
	}
	registry, err := coordinator.NewStaticRegistry(cfg.Worker.Coordinators)
	if err != nil {
		return nil, fmt.Errorf("failed to create static registry: %w", err)
	}
	logger.Info(ctx, "Using static coordinator discovery",
		slog.Any("coordinators", cfg.Worker.Coordinators))
	return coordinator.New(registry, clientConfig), nil
}
