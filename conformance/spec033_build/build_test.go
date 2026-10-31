package spec033_build_test

import (
	"os"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

func TestBuildWorkflowMaterializationLifecycle(t *testing.T) {
	flow := harness.NewRunner(t)
	flow.Mkdir(".home")
	flow.Mkdir(".xdg")
	env := []string{
		"FLOW_HOME=" + flow.ProjectPath(".flow"),
		"HOME=" + flow.ProjectPath(".home"),
		"XDG_CONFIG_HOME=" + flow.ProjectPath(".xdg"),
		"APPDATA=" + flow.ProjectPath(".home"),
		"USERPROFILE=" + flow.ProjectPath(".home"),
	}
	run := func(args ...string) *harness.Result {
		return flow.RunWithEnv(env, args...)
	}

	first := run("start", "build.yaml")
	first.ExpectExitCode(0)
	flow.ExpectTextFileContent("result.txt", "alpha\n")
	flow.ExpectTextFileContent("build-count.txt", "1\n")
	flow.ExpectTextFileContent("consume-count.txt", "1\n")
	flow.ExpectTextFileContent("controlled-count.txt", "1\n")
	flow.ExpectTextFileContent("always-count.txt", "1\n")

	second := run("start", "build.yaml")
	second.ExpectExitCode(0)
	require.Contains(t, second.Stdout(), "build: reuse (matched)")
	flow.ExpectTextFileContent("build-count.txt", "1\n")
	flow.ExpectTextFileContent("consume-count.txt", "1\n")
	flow.ExpectTextFileContent("controlled-count.txt", "1\n")
	flow.ExpectTextFileContent("always-count.txt", "2\n")

	dry := run("dry", "build.yaml")
	dry.ExpectExitCode(0)
	require.Contains(t, dry.Stdout(), "build: reuse (matched)")
	flow.ExpectTextFileContent("build-count.txt", "1\n")
	flow.ExpectTextFileContent("consume-count.txt", "1\n")

	forced := run("start", "--no-reuse", "build.yaml")
	forced.ExpectExitCode(0)
	require.Contains(t, forced.Stdout(), "build: execute (reuse_disabled)")
	flow.ExpectTextFileContent("build-count.txt", "2\n")
	flow.ExpectTextFileContent("consume-count.txt", "2\n")
	flow.ExpectTextFileContent("controlled-count.txt", "2\n")

	flow.WriteFile("source.txt", "beta\n")
	changed := run("start", "build.yaml")
	changed.ExpectExitCode(0)
	flow.ExpectTextFileContent("result.txt", "beta\n")
	flow.ExpectTextFileContent("build-count.txt", "3\n")
	flow.ExpectTextFileContent("consume-count.txt", "3\n")
	flow.ExpectTextFileContent("controlled-count.txt", "3\n")

	flow.WriteFile("fail.txt", "fail\n")
	failed := run("start", "--no-reuse", "build.yaml")
	failed.ExpectNonZeroExitCode()
	flow.ExpectTextFileContent("intermediate.txt", "beta\n")
	flow.ExpectTextFileContent("result.txt", "beta\n")

	require.NoError(t, os.Remove(flow.ProjectPath("fail.txt")))
	require.NoError(t, os.Remove(flow.ProjectPath("source.txt")))
	skipped := run("start", "--no-reuse", "build.yaml")
	skipped.ExpectExitCode(0)
	require.Contains(t, skipped.Stdout(), "build: none (precondition_not_met)")
	flow.ExpectTextFileContent("intermediate.txt", "beta\n")
	flow.ExpectTextFileContent("result.txt", "beta\n")
}

func TestBuildWorkflowRejectsDuplicateOutputProducers(t *testing.T) {
	flow := harness.NewRunner(t)
	result := flow.Run("dry", "duplicate-output.yaml")
	result.ExpectNonZeroExitCode()
	result.ExpectStderrContains("multiple producers", "shared.txt")
}

// A build step may publish a materialized path alongside typed values captured
// from its own output, so materialization must tolerate non-string entries.
func TestBuildStepPublishesTypedCapturedOutputs(t *testing.T) {
	flow := harness.NewRunner(t)
	flow.Mkdir(".home")
	flow.Mkdir(".xdg")
	env := []string{
		"FLOW_HOME=" + flow.ProjectPath(".flow"),
		"HOME=" + flow.ProjectPath(".home"),
		"XDG_CONFIG_HOME=" + flow.ProjectPath(".xdg"),
		"APPDATA=" + flow.ProjectPath(".home"),
		"USERPROFILE=" + flow.ProjectPath(".home"),
	}

	result := flow.RunWithEnv(env, "start", "typed-capture-output.yaml")
	result.ExpectExitCode(0)
	flow.ExpectTextFileContent("typed-capture.txt", "3\n")
}
