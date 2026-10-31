package cli_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

func TestPsShowsNothingWhenIdle(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	result := flow.RunWithEnv(env, "ps")
	result.ExpectExitCode(0)
	require.Contains(t, result.Stdout(), "No running processes")
}

func TestPsFiltersByDAGNameAndRunID(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-ps-filter-run"

	proc := startBackgroundAndWaitFor(t, flow, env, "started.out", "start", "--run-id="+runID, "long_running.yaml")
	defer proc.Stop()

	all := flow.RunWithEnv(env, "ps")
	all.ExpectExitCode(0)
	require.Contains(t, all.Stdout(), "long_running")
	require.Contains(t, all.Stdout(), runID)

	byName := flow.RunWithEnv(env, "ps", "-d", "long_running")
	byName.ExpectExitCode(0)
	require.Contains(t, byName.Stdout(), runID)

	wrongName := flow.RunWithEnv(env, "ps", "-d", "simple")
	wrongName.ExpectExitCode(0)
	require.Contains(t, wrongName.Stdout(), "No running processes")

	byRunID := flow.RunWithEnv(env, "ps", "-r", "ps-filter")
	byRunID.ExpectExitCode(0)
	require.Contains(t, byRunID.Stdout(), runID)

	wrongRunID := flow.RunWithEnv(env, "ps", "-r", "does-not-exist")
	wrongRunID.ExpectExitCode(0)
	require.Contains(t, wrongRunID.Stdout(), "No running processes")
}
