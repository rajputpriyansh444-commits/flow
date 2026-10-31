package dagrun_test

import (
	"testing"
	"time"

	"github.com/rajputpriyansh444-commits/flow/internal/dagrun"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/stretchr/testify/assert"
)

func TestIsLeaseActiveAllowsCoordinatorRestart(t *testing.T) {
	status := &ir.DAGRunStatus{LeaseAt: time.Now().Add(-time.Minute).UnixMilli()}

	assert.True(t, dagrun.IsLeaseActive(status, 0))
}
