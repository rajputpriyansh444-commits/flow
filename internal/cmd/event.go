package cmd

import (
	"context"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/logger"
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/logger/tag"
	"github.com/rajputpriyansh444-commits/flow/internal/eventstore"
	fileeventstore "github.com/rajputpriyansh444-commits/flow/internal/persis/file/eventstore"
)

func newEventService(ctx context.Context, cfg *config.Config) *eventstore.Service {
	if !cfg.EventStore.Enabled {
		return nil
	}
	store, err := fileeventstore.New(cfg.Paths.EventStoreDir)
	if err != nil {
		logger.Warn(ctx, "Failed to initialize event store; continuing without event persistence", tag.Error(err))
		return nil
	}
	return eventstore.New(store)
}
