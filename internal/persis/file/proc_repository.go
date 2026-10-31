package file

import (
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/rajputpriyansh444-commits/flow/internal/persis"
	procstore "github.com/rajputpriyansh444-commits/flow/internal/persis/file/proc"
)

// NewProcRepository connects the released .proc file layout to the shared repository.
func NewProcRepository(cfg *config.Config, opts ...procstore.StoreOption) *persis.ProcRepository {
	storeOpts := []procstore.StoreOption{
		procstore.WithStaleThreshold(cfg.Proc.StaleThreshold),
		procstore.WithHeartbeatInterval(cfg.Proc.HeartbeatInterval),
	}
	storeOpts = append(storeOpts, opts...)
	return persis.NewProcRepository(procstore.New(cfg.Paths.ProcDir, storeOpts...))
}
