package distr_test

import (
	"testing"
	"time"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/stretchr/testify/require"
)

const (
	distrTestProcHeartbeatInterval = 150 * time.Millisecond
	distrTestProcStaleThreshold    = 3 * time.Second
	distrTestProcSuccessTimeout    = 30 * time.Second
)

func TestExecution_ProcHeartbeat_DirectStart(t *testing.T) {
	f := newTestFixture(t, `
type: graph
name: worker-proc-heartbeat-start
worker_selector:
  test: "true"
steps:
  - name: sleep
    run: sleep 2
`, withProcConfig(distrTestProcHeartbeatInterval, distrTestProcStaleThreshold))
	defer f.cleanup()

	f.startScheduler(30 * time.Second)
	require.NoError(t, f.start())

	status := f.waitForStatus(ir.Succeeded, distrTestProcSuccessTimeout)
	require.Equal(t, ir.Succeeded, status.Status)
}

func TestExecution_ProcHeartbeat_QueuedDispatch(t *testing.T) {
	f := newTestFixture(t, `
type: graph
name: worker-proc-heartbeat-queued
worker_selector:
  test: "true"
steps:
  - name: sleep
    run: sleep 2
`, withProcConfig(distrTestProcHeartbeatInterval, distrTestProcStaleThreshold))
	defer f.cleanup()

	require.NoError(t, f.enqueue())
	f.waitForQueued()
	f.startScheduler(30 * time.Second)

	status := f.waitForStatus(ir.Succeeded, distrTestProcSuccessTimeout)
	require.Equal(t, ir.Succeeded, status.Status)
}
