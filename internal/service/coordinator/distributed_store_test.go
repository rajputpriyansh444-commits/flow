package coordinator

import (
	"path/filepath"

	"github.com/rajputpriyansh444-commits/flow/internal/persis/file"
	"github.com/rajputpriyansh444-commits/flow/internal/persis/store"
)

func newTestDispatchTaskStore(baseDir string) *store.DispatchTaskStore {
	return store.NewDispatchTaskStore(file.NewCollection(baseDir))
}

func newTestDispatchAdmissionTaskStore(baseDir string) (*store.DispatchTaskStore, *store.DAGRunLeaseStore, *store.ActiveDistributedRunStore) {
	leaseStore := newTestDAGRunLeaseStore(baseDir)
	activeStore := newTestActiveDistributedRunStore(baseDir)
	return store.NewDispatchTaskStore(
		file.NewCollection(baseDir),
		store.WithDispatchAdmissionLiveness(leaseStore, activeStore),
	), leaseStore, activeStore
}

func newTestWorkerHeartbeatStore(baseDir string) *store.WorkerHeartbeatStore {
	return store.NewWorkerHeartbeatStore(file.NewCollection(filepath.Join(baseDir, "workers")))
}

func newTestDAGRunLeaseStore(baseDir string) *store.DAGRunLeaseStore {
	return store.NewDAGRunLeaseStore(file.NewCollection(filepath.Join(baseDir, "leases")))
}

func newTestActiveDistributedRunStore(baseDir string) *store.ActiveDistributedRunStore {
	return store.NewActiveDistributedRunStore(file.NewCollection(filepath.Join(baseDir, "active-runs")))
}
