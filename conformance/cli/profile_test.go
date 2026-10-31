package cli_test

import (
	"regexp"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

// profileStatusPattern matches the "Status" row `profile show` prints, e.g.
// "Status       active".
var profileStatusPattern = regexp.MustCompile(`(?m)^Status\s+(\S+)`)

// profileStatus extracts the exact status value from `profile show` output,
// so callers can assert on the complete value instead of a substring that
// could also match a different status sharing part of the same word.
func profileStatus(stdout string) string {
	match := profileStatusPattern.FindStringSubmatch(stdout)
	if match == nil {
		return ""
	}
	return match[1]
}

func TestProfileCreateListShowDelete(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	create := flow.RunWithEnv(env, "profile", "create", "myprof", "--description=a test profile")
	create.ExpectExitCode(0)

	list := flow.RunWithEnv(env, "profile", "list")
	list.ExpectExitCode(0)
	require.Contains(t, list.Stdout(), "myprof")

	show := flow.RunWithEnv(env, "profile", "show", "myprof")
	show.ExpectExitCode(0)
	require.Contains(t, show.Stdout(), "a test profile")

	del := flow.RunWithEnv(env, "profile", "delete", "myprof")
	del.ExpectExitCode(0)

	listAfter := flow.RunWithEnv(env, "profile", "list")
	listAfter.ExpectExitCode(0)
	require.NotContains(t, listAfter.Stdout(), "myprof")
}

// TestProfileSetVarAppliesToRun proves a runtime profile variable actually
// flows into a run selecting it via --profile, exposed the same way step
// env values are: ${env.NAME}.
func TestProfileSetVarAppliesToRun(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	flow.RunWithEnv(env, "profile", "create", "myprof").ExpectExitCode(0)
	flow.RunWithEnv(env, "profile", "set-var", "myprof", "MY_VAR", "hello-from-profile").ExpectExitCode(0)

	show := flow.RunWithEnv(env, "profile", "show", "myprof")
	show.ExpectExitCode(0)
	require.Contains(t, show.Stdout(), "MY_VAR")
	require.Contains(t, show.Stdout(), "hello-from-profile")

	result := flow.RunWithEnv(env, "start", "--profile=myprof", "profile_consumer.yaml")
	result.ExpectExitCode(0)
	flow.ExpectFileContent("profile_var.out", "hello-from-profile\n")
}

func TestProfileEnableDisable(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	env := sharedEnv(t)

	flow.RunWithEnv(env, "profile", "create", "myprof").ExpectExitCode(0)

	disable := flow.RunWithEnv(env, "profile", "disable", "myprof")
	disable.ExpectExitCode(0)

	show := flow.RunWithEnv(env, "profile", "show", "myprof")
	show.ExpectExitCode(0)
	require.Equal(t, "disabled", profileStatus(show.Stdout()))

	enable := flow.RunWithEnv(env, "profile", "enable", "myprof")
	enable.ExpectExitCode(0)

	showAgain := flow.RunWithEnv(env, "profile", "show", "myprof")
	showAgain.ExpectExitCode(0)
	require.Equal(t, "active", profileStatus(showAgain.Stdout()))
}
