package cli_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

// TestCleanupIsAnAliasForRemoveHistory proves the deprecated `flow cleanup`
// command still works as documented: "Equivalent with rm: flow rm -H".
func TestCleanupIsAnAliasForRemoveHistory(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	flow.RunWithEnv(env, "start", "simple.yaml").ExpectExitCode(0)

	result := flow.RunWithEnv(env, "cleanup", "-y", "simple")
	result.ExpectExitCode(0)

	repeat := flow.RunWithEnv(env, "rm", "-H", "-f", "simple")
	repeat.ExpectExitCode(0)
	require.Contains(t, repeat.Stdout(), `No runs to delete for DAG "simple"`)
}

func TestCleanupDryRunPreviewsWithoutDeleting(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	flow.RunWithEnv(env, "start", "simple.yaml").ExpectExitCode(0)

	preview := flow.RunWithEnv(env, "cleanup", "--dry-run", "-y", "simple")
	preview.ExpectExitCode(0)
	require.Contains(t, preview.Stdout(), `Dry run: Would delete 1 run(s) for DAG "simple"`)

	result := flow.RunWithEnv(env, "rm", "-H", "-f", "simple")
	result.ExpectExitCode(0)
	require.Contains(t, result.Stdout(), `Successfully removed 1 run(s) for DAG "simple"`)
}
