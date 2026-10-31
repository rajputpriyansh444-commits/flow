package runtime

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/dispatch"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/stretchr/testify/assert"
)

func TestWorkerReportsClaim(t *testing.T) {
	t.Parallel()

	record := &dispatch.WorkerHeartbeatRecord{Stats: &dispatch.WorkerStats{
		RunningTasks: []*dispatch.RunningTask{{
			DAGRunID:   "parent-run",
			DAGName:    "parent",
			AttemptKey: "owner-key",
		}},
	}}
	childStatus := &ir.DAGRunStatus{Name: "child", DAGRunID: "child-run"}

	assert.True(t, workerReportsClaim(record, childStatus, "child-key", "owner-key"))
	assert.False(t, workerReportsClaim(record, childStatus, "child-key", "different-claim"))
}
