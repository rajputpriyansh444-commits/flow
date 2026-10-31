package spec019_sub_dag_working_directory_test

import (
	"runtime"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

func TestSubDAGWorkingDirectoryRuntimeUnix(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("fixtures use POSIX shell snippets")
	}

	t.Run("parent working directory is not inherited by default", func(t *testing.T) {
		t.Parallel()

		flow := harness.NewRunner(t)
		result := flow.Run("start", "parent_workdir_not_inherited.yaml")
		result.ExpectExitCode(0)
	})

	t.Run("parallel children use isolated default work directories", func(t *testing.T) {
		t.Parallel()

		flow := harness.NewRunner(t)
		result := flow.Run("start", "parallel_children_default_workdir.yaml")
		result.ExpectExitCode(0)
		flow.ExpectFileContains(
			"parallel-workdir-results.txt",
			`"total": 2`,
			`"succeeded": 2`,
			`"CHILD_INFO"`,
			"item=alpha",
			"item=beta",
		)
		flow.ExpectNoFile("collision.txt")
	})

	t.Run("explicit child working directory wins", func(t *testing.T) {
		t.Parallel()

		flow := harness.NewRunner(t)
		result := flow.Run("start", "explicit_child_working_dir.yaml")
		result.ExpectExitCode(0)
		flow.ExpectFileContains(
			"explicit-child-results.txt",
			`"CHILD_INFO"`,
			"cwd=shared-child-work",
			"comparison=different",
			"context=same",
		)
	})

	t.Run("base config child working directory is explicit", func(t *testing.T) {
		t.Parallel()

		flow := harness.NewRunner(t)
		flow.WriteFile("base.yaml", "working_dir: base-child-work\n")
		result := flow.RunWithEnv(
			[]string{"FLOW_BASE_CONFIG=" + flow.ProjectPath("base.yaml")},
			"start",
			"base_config_child_working_dir.yaml",
		)
		result.ExpectExitCode(0)
		flow.ExpectFileContains(
			"base-config-child-results.txt",
			`"CHILD_INFO"`,
			"cwd=base-child-work",
			"comparison=different",
			"context=same",
		)
	})
}
