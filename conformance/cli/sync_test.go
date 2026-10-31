package cli_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

func TestSyncDisabled(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	result := flow.RunWithEnv(sharedEnv(t), "sync", "status")
	result.ExpectNonZeroExitCode()
	result.ExpectStderrContains("git sync is not enabled")
}

// Status reports configuration without contacting the remote repository.
func TestSyncConfiguredStatus(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	const repository = "https://example.invalid/flow/workflows.git"
	env := append(sharedEnv(t),
		"FLOW_GITSYNC_ENABLED=true",
		"FLOW_GITSYNC_REPOSITORY="+repository,
		"FLOW_GITSYNC_BRANCH=main",
		"FLOW_GITSYNC_AUTH_TOKEN=conformance-token",
	)

	result := flow.RunWithEnv(env, "sync", "status")
	result.ExpectExitCode(0)
	requireLineContains(t, result.Stdout(), "Repository:", repository)
	requireLineContains(t, result.Stdout(), "Branch:", "main")
}
