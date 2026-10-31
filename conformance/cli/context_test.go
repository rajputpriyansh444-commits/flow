package cli_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

func TestContextListShowsBuiltInLocal(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	result := flow.RunWithEnv(env, "context", "list")
	result.ExpectExitCode(0)
	require.Contains(t, result.Stdout(), "local")
}

func TestContextAddUpdateRemove(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	add := flow.RunWithEnv(env, "context", "add", "myremote",
		"--server=http://127.0.0.1:1", "--api-key=flow_testkey123456789", "--description=first")
	add.ExpectExitCode(0)

	list := flow.RunWithEnv(env, "context", "list")
	list.ExpectExitCode(0)
	require.Contains(t, list.Stdout(), "myremote")
	require.Contains(t, list.Stdout(), "first")

	update := flow.RunWithEnv(env, "context", "update", "myremote", "--description=second")
	update.ExpectExitCode(0)

	listAfterUpdate := flow.RunWithEnv(env, "context", "list")
	listAfterUpdate.ExpectExitCode(0)
	require.Contains(t, listAfterUpdate.Stdout(), "second")
	require.NotContains(t, listAfterUpdate.Stdout(), "first")

	remove := flow.RunWithEnv(env, "context", "remove", "myremote")
	remove.ExpectExitCode(0)

	listAfterRemove := flow.RunWithEnv(env, "context", "list")
	listAfterRemove.ExpectExitCode(0)
	require.NotContains(t, listAfterRemove.Stdout(), "myremote")
}

func TestContextUseSwitchesCurrent(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	flow.RunWithEnv(env, "context", "add", "myremote", "--server=http://127.0.0.1:1", "--api-key=flow_testkey123456789").ExpectExitCode(0)

	use := flow.RunWithEnv(env, "context", "use", "myremote")
	use.ExpectExitCode(0)

	list := flow.RunWithEnv(env, "context", "list")
	list.ExpectExitCode(0)
	require.Regexp(t, `(?m)^\*\s+myremote\b`, list.Stdout())

	backToLocal := flow.RunWithEnv(env, "context", "use", "local")
	backToLocal.ExpectExitCode(0)

	listAfter := flow.RunWithEnv(env, "context", "list")
	listAfter.ExpectExitCode(0)
	require.Regexp(t, `(?m)^\*\s+local\b`, listAfter.Stdout())
}

func TestContextTestLocalIsAlwaysAvailable(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	result := flow.RunWithEnv(env, "context", "test", "local")
	result.ExpectExitCode(0)
	require.Contains(t, result.Stdout(), "always available")
}

func TestContextTestUnreachableRemoteFails(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	flow.RunWithEnv(env, "context", "add", "myremote", "--server=http://127.0.0.1:1", "--api-key=flow_testkey123456789").ExpectExitCode(0)

	result := flow.RunWithEnv(env, "context", "test", "myremote")
	result.ExpectNonZeroExitCode()
}
