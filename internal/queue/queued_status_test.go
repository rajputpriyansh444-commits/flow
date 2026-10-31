package queue_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/queue"
	"github.com/stretchr/testify/require"
)

func TestPreservedQueueTriggerType(t *testing.T) {
	t.Parallel()

	require.Equal(t, ir.TriggerTypeWebhook, queue.PreservedQueueTriggerType(&ir.DAGRunStatus{
		Status:      ir.Queued,
		TriggerType: ir.TriggerTypeWebhook,
	}))
	require.Equal(t, ir.TriggerTypeCatchUp, queue.PreservedQueueTriggerType(&ir.DAGRunStatus{
		Status:      ir.Queued,
		TriggerType: ir.TriggerTypeCatchUp,
	}))
	require.Equal(t, ir.TriggerTypeUnknown, queue.PreservedQueueTriggerType(&ir.DAGRunStatus{
		Status:      ir.Queued,
		TriggerType: ir.TriggerTypeRetry,
	}))
	require.Equal(t, ir.TriggerTypeUnknown, queue.PreservedQueueTriggerType(&ir.DAGRunStatus{
		Status:      ir.Succeeded,
		TriggerType: ir.TriggerTypeWebhook,
	}))
}
