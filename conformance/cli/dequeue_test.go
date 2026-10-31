package cli_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

// Dequeue removes the oldest run and leaves the successor queued.
func TestDequeueFirstFromQueue(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-dequeue-first"
	const nextRunID = "cli-dequeue-next"

	flow.RunWithEnv(env, "enqueue", "--run-id="+runID, "simple.yaml").ExpectExitCode(0)
	flow.RunWithEnv(env, "enqueue", "--run-id="+nextRunID, "simple.yaml").ExpectExitCode(0)

	result := flow.RunWithEnv(env, "dequeue", "simple")
	result.ExpectExitCode(0)

	status := flow.RunWithEnv(env, "status", "--run-id="+runID, "simple.yaml")
	status.ExpectNonZeroExitCode()
	status.ExpectStderrContains("dag-run ID not found")

	next := flow.RunWithEnv(env, "status", "--run-id="+nextRunID, "simple.yaml")
	next.ExpectExitCode(0)
	require.Equal(t, "Queued", resultStatus(next.Stdout()))
}

// TestDequeueSpecificDAGRun proves `--dag-run` selects exactly one queued
// item to remove, leaving the rest of the queue untouched.
func TestDequeueSpecificDAGRun(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const keepRunID = "cli-dequeue-keep"
	const removeRunID = "cli-dequeue-remove"

	flow.RunWithEnv(env, "enqueue", "--run-id="+keepRunID, "simple.yaml").ExpectExitCode(0)
	flow.RunWithEnv(env, "enqueue", "--run-id="+removeRunID, "simple.yaml").ExpectExitCode(0)

	result := flow.RunWithEnv(env, "dequeue", "simple", "--dag-run=simple:"+removeRunID)
	result.ExpectExitCode(0)

	removedStatus := flow.RunWithEnv(env, "status", "--run-id="+removeRunID, "simple.yaml")
	removedStatus.ExpectNonZeroExitCode()
	removedStatus.ExpectStderrContains("dag-run ID not found")

	keptStatus := flow.RunWithEnv(env, "status", "--run-id="+keepRunID, "simple.yaml")
	keptStatus.ExpectExitCode(0)
	require.Contains(t, keptStatus.Stdout(), "Queued")
}

func TestDequeueEmptyQueueFails(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	result := flow.RunWithEnv(env, "dequeue", "simple")
	result.ExpectNonZeroExitCode()
	result.ExpectStderrContains("no dag-run found in queue")
}
