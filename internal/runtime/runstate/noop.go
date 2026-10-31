package runstate

import (
	"context"

	"github.com/rajputpriyansh444-commits/flow/internal/dagrun"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
)

// NewNoopAttempt creates execution state for a run persisted outside this process.
func NewNoopAttempt(req BeginAttemptRequest) Attempt {
	id := req.AttemptID
	if id == "" {
		id = req.RunID
	}
	return noopAttempt{Attempt: dagrun.NewNoopAttempt(id, req.DAG)}
}

type noopAttempt struct {
	dagrun.Attempt
}

func (a noopAttempt) RecordStatus(ctx context.Context, status ir.DAGRunStatus) error {
	return a.Write(ctx, status)
}

func (a noopAttempt) RecordOutputs(ctx context.Context, outputs *ir.DAGRunOutputs) error {
	return a.WriteOutputs(ctx, outputs)
}

func (a noopAttempt) RequestCancel(ctx context.Context) error {
	return a.Abort(ctx)
}

func (a noopAttempt) CancelRequested(ctx context.Context) (bool, error) {
	return a.IsAborting(ctx)
}

func (noopAttempt) MaterializeWorkDir(context.Context) (string, error) {
	return "", nil
}

func (noopAttempt) SnapshotWorkDir(context.Context, string) error {
	return nil
}
