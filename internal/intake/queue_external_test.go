package intake_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/rajputpriyansh444-commits/flow/internal/intake"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/persis"
	"github.com/rajputpriyansh444-commits/flow/internal/persis/file"
	"github.com/rajputpriyansh444-commits/flow/internal/persis/store"
	"github.com/rajputpriyansh444-commits/flow/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestEnqueueRunDoesNotSeedQueuedCondition(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tmp := t.TempDir()
	dag := &ir.DAG{Name: "queued-condition"}
	ir.InitializeDefaults(dag)
	dagRunRepository := testutil.NewFileDAGRunRepository(filepath.Join(tmp, "dag-runs"), persis.DAGRunRepositoryOptions{})
	queueStore := store.NewQueueStore(file.NewCollection(filepath.Join(tmp, "queue")))
	now := time.Date(2026, 5, 19, 1, 2, 3, 0, time.UTC)

	_, err := intake.EnqueueRun(ctx, intake.QueueRequest{
		DAGRunRepository: dagRunRepository,
		QueueStore:       queueStore,
		DAG:              dag,
		DAGRunID:         "run-1",
		LogBaseDir:       filepath.Join(tmp, "logs"),
		Now:              func() time.Time { return now },
	})
	require.NoError(t, err)

	attempt, err := dagRunRepository.FindAttempt(ctx, ir.NewDAGRunRef(dag.Name, "run-1"))
	require.NoError(t, err)
	status, err := attempt.ReadStatus(ctx)
	require.NoError(t, err)
	require.Equal(t, ir.Queued, status.Status)
	require.Empty(t, status.Conditions)
}
