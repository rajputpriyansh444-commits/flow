package cmd

import (
	"context"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
)

func RestoreDAGFromStatusForTest(ctx context.Context, dag *ir.DAG, status *ir.DAGRunStatus) (*ir.DAG, error) {
	return restoreDAGFromStatus(ctx, dag, status, nil)
}

var WorkerFlagsForTest = workerFlags
