package cli_test

import (
	"runtime"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

// TestStopTerminatesActiveRun proves the primary `flow stop --run-id=<id>
// <dag>` contract: it gracefully terminates an active run, which then
// reports a terminal, non-succeeded status rather than continuing to run to
// completion. Uses the same deterministic poll-then-stop pattern as
// conformance/spec018_parallel_foreach's abort test rather than racing a
// wall-clock delay against process-launch overhead.
//
// The command's second documented mode -- cancelling a root run that failed
// and is still pending its DAG-level auto-retry -- is not covered here; it
// needs a DAG-level retry_policy fixture and is left as a follow-up.
func TestStopTerminatesActiveRun(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-stop-active"

	proc := startBackgroundAndWaitFor(t, flow, env, "started.out", "start", "--run-id="+runID, "long_running.yaml")
	defer proc.Stop()

	result := flow.RunWithEnv(env, "stop", "--run-id="+runID, "long_running.yaml")
	result.ExpectExitCode(0)

	waitForProcessDone(t, proc)

	// Assert the specific documented terminal status rather than only
	// excluding "Succeeded"/"Running": those negative checks would also
	// pass for a stalled non-terminal status such as "Queued" if status
	// finalization got stuck after the process exited.
	waitForStatus(t, flow, env, runID, "long_running.yaml", "Aborted")
}

// Step cleanup must leave abort and exit handlers time to complete.
func TestStopPreservesHandlers(t *testing.T) {
	t.Parallel()
	files := []string{"stop_handlers_default.yaml", "stop_handlers_forced.yaml"}
	if runtime.GOOS == "windows" {
		files = []string{"stop_handlers_default_windows.yaml", "stop_handlers_forced_windows.yaml"}
	}
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			t.Parallel()
			flow := harness.NewRunner(t)
			env := append(sharedEnv(t), "FLOW_SIGNAL_PROPAGATION=false")
			const runID = "cli-stop-handlers"
			proc := startBackgroundAndWaitFor(t, flow, env, "started.out", "start", "--run-id="+runID, file)
			defer proc.Stop()

			flow.RunWithEnv(env, "stop", "--run-id="+runID, file).ExpectExitCode(0)
			waitForProcessDone(t, proc)
			flow.ExpectTextFileContent("handlers.out", "abort\nexit\n")
			waitForStatus(t, flow, env, runID, file, "Aborted")
		})
	}
}
