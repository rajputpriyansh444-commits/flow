package persis

import (
	"context"
	"errors"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/proc"
)

// ProcStore persists process heartbeats for [ProcRepository].
type ProcStore interface {
	Validate(ctx context.Context) error
	WithLock(ctx context.Context, groupName string, fn func() error) error
	Acquire(ctx context.Context, groupName string, meta proc.ProcMeta) (proc.ProcHandle, error)
	ListEntries(ctx context.Context, groupName string) ([]proc.ProcEntry, error)
	LatestHeartbeat(ctx context.Context, groupName string, dagRun ir.DAGRunRef) (*proc.ProcHeartbeat, error)
	ListAllEntries(ctx context.Context) ([]proc.ProcEntry, error)
	RemoveIfStale(ctx context.Context, entry proc.ProcEntry) error
}

// ProcLockError reports a failure to acquire a process-group lock.
type ProcLockError struct {
	err error
}

// NewProcLockError classifies a process-group lock-acquisition failure.
func NewProcLockError(err error) error {
	if err == nil {
		return nil
	}
	return &ProcLockError{err: err}
}

func (e *ProcLockError) Error() string { return e.err.Error() }
func (e *ProcLockError) Unwrap() error { return e.err }

// IsProcLockError reports whether err represents process-group lock acquisition failure.
func IsProcLockError(err error) bool {
	var target *ProcLockError
	return errors.As(err, &target)
}
