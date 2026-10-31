package cli_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

// TestRetryReexecutesStep proves the basic contract: retrying a previously
// executed run with the same run ID actually re-runs the step, not just
// re-reports the old status. The step appends a marker each time it runs, so
// two attempts leave two marker lines.
func TestRetryReexecutesStep(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-retry-run"

	first := flow.RunWithEnv(env, "start", "--run-id="+runID, "failing.yaml")
	first.ExpectNonZeroExitCode()
	flow.ExpectFileContent("attempts.out", "attempt\n")

	retry := flow.RunWithEnv(env, "retry", "--run-id="+runID, "failing")
	retry.ExpectNonZeroExitCode()
	flow.ExpectFileContent("attempts.out", "attempt\nattempt\n")
}

// TestRetryPartialForeach proves retry recovers a partial foreach and keeps
// unrelated successful work complete.
func TestRetryPartialForeach(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-retry-partial-foreach"

	first := flow.RunWithEnv(env, "start", "--run-id="+runID, "retry_partial_foreach.yaml")
	first.ExpectExitCode(0)
	status := flow.RunWithEnv(env, "status", "--run-id="+runID, "retry_partial_foreach")
	status.ExpectExitCode(0)
	require.Contains(t, status.Stdout(), "Result: Partially Succeeded")
	flow.ExpectFileContent("foreach-a.out", "run\n")
	flow.ExpectFileContent("foreach-b.out", "run\n")
	flow.ExpectFileContent("independent.out", "independent\n")
	flow.ExpectFileContent("foreach-downstream.out", "downstream\n")

	flow.WriteFile("allow-b.txt", "")
	retry := flow.RunWithEnv(env, "retry", "--run-id="+runID, "retry_partial_foreach")
	retry.ExpectExitCode(0)
	status = flow.RunWithEnv(env, "status", "--run-id="+runID, "retry_partial_foreach")
	status.ExpectExitCode(0)
	require.Contains(t, status.Stdout(), "Result: Succeeded")
	flow.ExpectFileContent("foreach-a.out", "run\nrun\n")
	flow.ExpectFileContent("foreach-b.out", "run\nrun\n")
	flow.ExpectFileContent("independent.out", "independent\n")
	flow.ExpectFileContent("foreach-downstream.out", "downstream\ndownstream\n")
}

// TestRetryPartialParallel proves retry recovers a partial child DAG and keeps
// unrelated successful work complete.
func TestRetryPartialParallel(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-retry-partial-parallel"

	first := flow.RunWithEnv(env, "start", "--run-id="+runID, "retry_partial_parallel.yaml")
	first.ExpectExitCode(0)
	status := flow.RunWithEnv(env, "status", "--run-id="+runID, "retry_partial_parallel")
	status.ExpectExitCode(0)
	require.Contains(t, status.Stdout(), "Result: Partially Succeeded")
	flow.ExpectFileContent("parallel-child-runs.out", "child\n")
	flow.ExpectFileContent("parallel-stable-runs.out", "stable\n")
	flow.ExpectFileContent("independent.out", "independent\n")
	flow.ExpectFileContent("parallel-downstream.out", "downstream\n")

	flow.WriteFile("allow-child", "")
	retry := flow.RunWithEnv(env, "retry", "--run-id="+runID, "retry_partial_parallel")
	retry.ExpectExitCode(0)
	status = flow.RunWithEnv(env, "status", "--run-id="+runID, "retry_partial_parallel")
	status.ExpectExitCode(0)
	require.Contains(t, status.Stdout(), "Result: Succeeded")
	flow.ExpectFileContent("parallel-child-runs.out", "child\nchild\n")
	flow.ExpectFileContent("parallel-stable-runs.out", "stable\n")
	flow.ExpectFileContent("independent.out", "independent\n")
	flow.ExpectFileContent("parallel-downstream.out", "downstream\ndownstream\n")
}

// TestRetryStepFlagReexecutesOnlyThatStep proves --step scopes the retry to
// one step: an earlier, already-succeeded step must not re-run.
func TestRetryStepFlagReexecutesOnlyThatStep(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-retry-step"

	first := flow.RunWithEnv(env, "start", "--run-id="+runID, "multi_step_failing.yaml")
	first.ExpectNonZeroExitCode()
	flow.ExpectFileContent("steps.out", "first\nsecond\n")

	retry := flow.RunWithEnv(env, "retry", "--run-id="+runID, "--step=second", "multi_step_failing")
	retry.ExpectNonZeroExitCode()
	flow.ExpectFileContent("steps.out", "first\nsecond\nsecond\n")
}

// TestRetryReplacesStdoutArtifact proves a stdout artifact holds the output of
// the latest attempt, not the output of every attempt of the run.
func TestRetryReplacesStdoutArtifact(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-retry-artifact"
	const artifact = "artifacts/*/*/*/*_artifact_retry_*/report.txt"

	first := flow.RunWithEnv(env, "start", "--run-id="+runID, "artifact_retry.yaml")
	first.ExpectNonZeroExitCode()
	flow.ExpectGlobFileContent(artifact, "line\n")

	flow.WriteFile("ok", "")
	retry := flow.RunWithEnv(env, "retry", "--run-id="+runID, "artifact_retry")
	retry.ExpectExitCode(0)
	flow.ExpectGlobFileContent(artifact, "line\n")
}

func TestRetryDownstreamRequiresStep(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)
	const runID = "cli-retry-downstream"

	first := flow.RunWithEnv(env, "start", "--run-id="+runID, "failing.yaml")
	first.ExpectNonZeroExitCode()

	retry := flow.RunWithEnv(env, "retry", "--run-id="+runID, "--downstream", "failing")
	retry.ExpectNonZeroExitCode()
	retry.ExpectStderrContains("--downstream requires --step")
}

func TestRetryMissingRunIDFails(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	result := flow.RunWithEnv(env, "retry", "failing")
	result.ExpectNonZeroExitCode()
	result.ExpectStderrContains("run-id")
}
