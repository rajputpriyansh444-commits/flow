// This file covers flag-level depth for commands whose core behavior is
// already exercised elsewhere: enqueue, status, history, dry, and start.
package cli_test

import (
	"regexp"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

var subdagIDPattern = regexp.MustCompile(`subdag:\s*(\S+)`)

// TestStatusSubRunIDAddressesChildRun proves --sub-run-id lets `flow status`
// address a child DAG run directly, and that it requires --run-id.
func TestStatusSubRunIDAddressesChildRun(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-depth-parent"

	started := flow.RunWithEnv(env, "start", "--run-id="+runID, "parent_child.yaml")
	started.ExpectExitCode(0)

	match := subdagIDPattern.FindStringSubmatch(started.Stdout())
	require.Lenf(t, match, 2, "expected a subdag id in start output:\n%s", started.Stdout())
	childRunID := match[1]

	withoutRunID := flow.RunWithEnv(env, "status", "--sub-run-id="+childRunID, "parent_child.yaml")
	withoutRunID.ExpectNonZeroExitCode()
	withoutRunID.ExpectStderrContains("--sub-run-id requires --run-id")

	withRunID := flow.RunWithEnv(env, "status", "--run-id="+runID, "--sub-run-id="+childRunID, "parent_child.yaml")
	withRunID.ExpectExitCode(0)
	require.Contains(t, withRunID.Stdout(), "dag: depth_child")
}

// TestEnqueueQueueFlagSelectsQueue proves --queue overrides the DAG's
// default process group as the queue an enqueued run lands in.
func TestEnqueueQueueFlagSelectsQueue(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-depth-enqueue-queue"

	flow.RunWithEnv(env, "enqueue", "--run-id="+runID, "--queue=custom-queue", "simple.yaml").ExpectExitCode(0)

	// The default queue (the DAG's own name) must be empty...
	wrongQueue := flow.RunWithEnv(env, "dequeue", "simple")
	wrongQueue.ExpectNonZeroExitCode()
	wrongQueue.ExpectStderrContains("no dag-run found in queue")

	// ...while the overridden queue holds the item.
	rightQueue := flow.RunWithEnv(env, "dequeue", "custom-queue")
	rightQueue.ExpectExitCode(0)
}

// TestEnqueueLabelsFlowIntoHistory proves --labels on enqueue is recorded and
// that `flow history --labels` filters by it.
func TestEnqueueLabelsFlowIntoHistory(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const labeledRunID = "cli-depth-labeled"
	const plainRunID = "cli-depth-plain"

	flow.RunWithEnv(env, "enqueue", "--run-id="+labeledRunID, "--labels=env=prod", "simple.yaml").ExpectExitCode(0)
	flow.RunWithEnv(env, "enqueue", "--run-id="+plainRunID, "simple.yaml").ExpectExitCode(0)

	result := flow.RunWithEnv(env, "history", "simple", "--labels=env=prod")
	result.ExpectExitCode(0)
	require.Contains(t, result.Stdout(), labeledRunID)
	require.NotContains(t, result.Stdout(), plainRunID)
}

// TestHistoryStatusFilterAndLimit proves --status filters by outcome and
// --limit caps how many results are returned.
func TestHistoryStatusFilterAndLimit(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	flow.RunWithEnv(env, "start", "--run-id=cli-depth-hist-ok-1", "simple.yaml").ExpectExitCode(0)
	flow.RunWithEnv(env, "start", "--run-id=cli-depth-hist-ok-2", "simple.yaml").ExpectExitCode(0)
	flow.RunWithEnv(env, "start", "--run-id=cli-depth-hist-fail-1", "failing.yaml")

	onlyFailed := flow.RunWithEnv(env, "history", "--status=failed")
	onlyFailed.ExpectExitCode(0)
	require.Contains(t, onlyFailed.Stdout(), "cli-depth-hist-fail-1")
	require.NotContains(t, onlyFailed.Stdout(), "cli-depth-hist-ok-1")
	require.NotContains(t, onlyFailed.Stdout(), "cli-depth-hist-ok-2")

	limited := flow.RunWithEnv(env, "history", "simple", "--limit=1")
	limited.ExpectExitCode(0)
	matches := regexp.MustCompile(`cli-depth-hist-ok-\d`).FindAllString(limited.Stdout(), -1)
	require.Len(t, matches, 1, "expected --limit=1 to cap results to one run:\n%s", limited.Stdout())
}

// TestDryRunResolvesParamsWithoutExecuting proves `flow dry --params`
// resolves parameter values (a missing required param would fail dry too)
// while still performing no real action.
func TestDryRunResolvesParamsWithoutExecuting(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	result := flow.RunWithEnv(env, "dry", "--params", "greeting=hello-dry", "params_dag.yaml")
	result.ExpectExitCode(0)
	flow.ExpectNoFile("params_dag.out")

	missing := flow.RunWithEnv(env, "dry", "params_dag.yaml")
	missing.ExpectNonZeroExitCode()
}

// TestStartLabelsFlowIntoHistory proves --labels on `start` is recorded the
// same way it is for enqueue.
func TestStartLabelsFlowIntoHistory(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-depth-start-labeled"

	flow.RunWithEnv(env, "start", "--run-id="+runID, "--labels=env=prod", "simple.yaml").ExpectExitCode(0)

	result := flow.RunWithEnv(env, "history", "simple", "--labels=env=prod")
	result.ExpectExitCode(0)
	require.Contains(t, result.Stdout(), runID)
}
