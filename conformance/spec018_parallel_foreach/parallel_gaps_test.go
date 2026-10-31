package spec018_parallel_foreach_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

// Enqueue accepts every child without a scheduler executing the runs.
func TestParallelEnqueue(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	result := flow.Run("start", "parallel_dag_enqueue.yaml")
	result.ExpectExitCode(0)
	flow.ExpectFileContains(
		"parallel-enqueue-results.txt",
		`"total": 2`,
		`"queued": 2`,
		`"name": "child-enqueue-item"`,
		`"queue": "parallel-enqueue-queue"`,
	)
}

// Stop must report an aborted run after the first fan-out item starts.
func TestParallelAbort(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := []string{"FLOW_HOME=" + filepath.Join(t.TempDir(), "flow")}
	const runID = "spec018-abort"

	proc := flow.StartWithEnv(env, "start", "--run-id="+runID, "parallel_timeout_abort.yaml")

	deadline := time.Now().Add(harness.WaitTimeout(t))
	for {
		// Redirection creates the file before printf writes the marker.
		content, err := os.ReadFile(flow.ProjectPath("started-one.txt"))
		if err == nil && string(content) == "started\n" {
			break
		}
		if err != nil && !os.IsNotExist(err) {
			t.Fatalf("reading start marker: %v", err)
		}
		if time.Now().After(deadline) {
			t.Fatalf("item one never started: %s", proc.FailureOutput())
		}
		select {
		case <-proc.Done():
			t.Fatalf("flow start exited before item one started: %s", proc.FailureOutput())
		case <-time.After(50 * time.Millisecond):
		}
	}
	stopResult := flow.RunWithEnv(env, "stop", "--run-id="+runID, "parallel_timeout_abort.yaml")
	stopResult.ExpectExitCode(0)

	select {
	case <-proc.Done():
	case <-time.After(harness.WaitTimeout(t)):
		t.Fatal("flow start did not exit after flow stop returned")
	}

	status := flow.RunWithEnv(env, "status", "--run-id="+runID, "parallel_timeout_abort.yaml")
	status.ExpectExitCode(0)
	require.Contains(t, status.Stdout(), "Aborted")
}

// A tolerated child failure propagates partially_succeeded to the parent.
func TestParallelPartial(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := []string{"FLOW_HOME=" + filepath.Join(t.TempDir(), "flow")}
	const runID = "spec018-partial"
	result := flow.RunWithEnv(env, "start", "--run-id="+runID, "parallel_partially_succeeded.yaml")
	result.ExpectExitCode(0)
	flow.ExpectFileContains(
		"parallel-partial-results.txt",
		`"total": 2`,
		`"succeeded": 2`,
		`"failed": 0`,
	)

	status := flow.RunWithEnv(env, "status", "--run-id="+runID, "parallel_partially_succeeded.yaml")
	status.ExpectExitCode(0)
	require.Contains(t, status.Stdout(), "Partially Succeeded")
}

// The step outputs channel carries one entry per successful child DAG run in
// parallel item order, and a failed child contributes no entry, so index N
// addresses the Nth successful child rather than the Nth item.
func TestParallelOutputsChannel(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	// The failing child makes the fan-out report an error, which the parent
	// tolerates through continue_on but still surfaces as a non-zero exit.
	result := flow.Run("start", "parallel_outputs_channel.yaml")
	result.ExpectExitCode(1)

	flow.ExpectFileContains(
		"parallel-outputs-channel.txt",
		`[{"RESULT":"alpha"},{"RESULT":"gamma"}]`,
	)
	flow.ExpectFileContains(
		"parallel-outputs-refs.txt",
		"first=alpha",
		"second=gamma",
	)
}
