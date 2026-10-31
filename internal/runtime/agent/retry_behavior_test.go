package agent

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/stretchr/testify/assert"
)

func TestAgentCurrentAutoRetryCount(t *testing.T) {
	t.Parallel()

	t.Run("InitialAttemptStartsAtZero", func(t *testing.T) {
		t.Parallel()
		a := &Agent{}
		assert.Equal(t, 0, a.currentAutoRetryCount())
	})

	t.Run("RetryTargetPreservesPersistedCount", func(t *testing.T) {
		t.Parallel()
		a := &Agent{
			retryTarget: &ir.DAGRunStatus{AutoRetryCount: 1},
		}
		assert.Equal(t, 1, a.currentAutoRetryCount())
	})
}
