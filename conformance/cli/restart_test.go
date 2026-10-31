package cli_test

import (
	"regexp"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

// restartedRunIDPattern matches restart's own announcement of the new
// dag-run it starts in place of the one it stopped, e.g.:
// `msg="Dag-run restart initiated" dag=restartable run-id=<new-id> file=...`
var restartedRunIDPattern = regexp.MustCompile(`Dag-run restart initiated" dag=\S+ run-id=(\S+)`)

// Restart aborts the original run and completes a distinct replacement.
func TestRestartReplacesRun(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const originalRunID = "cli-restart-original"

	proc := startBackgroundAndWaitFor(t, flow, env, "restart_marks.out", "start", "--run-id="+originalRunID, "restartable.yaml")
	defer proc.Stop()

	restart := flow.RunWithEnv(env, "restart", "restartable")
	restart.ExpectExitCode(0)

	match := restartedRunIDPattern.FindStringSubmatch(restart.Stderr())
	require.Lenf(t, match, 2, "expected restart to announce the new run's id:\n%s", restart.Stderr())
	require.NotEqual(t, originalRunID, match[1], "restart must start a run with a new id, not reuse the original")

	waitForProcessDone(t, proc)

	flow.ExpectFileContent("restart_marks.out", "started\nstarted\n")

	// Finalizing the aborted original attempt is asynchronous with the
	// process actually exiting, so poll for its terminal status rather than
	// checking once immediately after the process exits.
	waitForStatus(t, flow, env, originalRunID, "restartable.yaml", "Aborted")
	waitForStatus(t, flow, env, match[1], "restartable.yaml", "Succeeded")
}
