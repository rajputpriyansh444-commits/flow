package spec064_dry_run_step_checks_test

import (
	"runtime"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

const dryRunWarning = "Dry run: step may fail on this host"

func TestDryAcceptsCommandAndShell(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		file   string
		output string
	}{
		{name: "command resolvable on PATH", file: "valid_command.yaml", output: "valid-command.out"},
		{name: "step-level shell resolvable on PATH", file: "valid_shell.yaml", output: "valid-shell.out"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("dry", tc.file)
			result.ExpectExitCode(0)
			// A resolvable command or shell produces no warning.
			result.ExpectStderrNotContains(dryRunWarning)
			// Dry validation must not create step output files.
			flow.ExpectNoFile(tc.output)
		})
	}
}

// A missing executable is a warning, not a failure: the real run may execute
// on another host, or an upstream step may install the executable first.
func TestDryWarnsMissingExecutable(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		file    string
		missing string
		output  string
	}{
		{
			name:    "command",
			file:    "missing_command.yaml",
			missing: "flow-conformance-missing-command-9f3c2b1a",
			output:  "missing-command.out",
		},
		{
			name:    "shell",
			file:    "missing_shell.yaml",
			missing: "flow-conformance-missing-shell-9f3c2b1a",
			output:  "missing-shell.out",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("dry", tc.file)
			result.ExpectExitCode(0)
			result.ExpectStderrContains(dryRunWarning, tc.missing)
			flow.ExpectNoFile(tc.output)
		})
	}
}

func TestDryChecksExecutePermission(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("execute permission bits are not meaningful on Windows")
	}
	t.Parallel()

	const script = "#!/bin/sh\nprintf 'ran\\n' > script-ran.out\n"

	t.Run("executable script is accepted", func(t *testing.T) {
		t.Parallel()

		flow := harness.NewRunner(t)
		flow.WriteExecutable("scripts/step.sh", script)

		result := flow.Run("dry", "script_permission.yaml")
		result.ExpectExitCode(0)
		result.ExpectStderrNotContains(dryRunWarning)
		flow.ExpectNoFile("script-ran.out")
	})

	t.Run("non-executable script warns", func(t *testing.T) {
		t.Parallel()

		flow := harness.NewRunner(t)
		flow.WriteFile("scripts/step.sh", script)

		result := flow.Run("dry", "script_permission.yaml")
		result.ExpectExitCode(0)
		result.ExpectStderrContains(dryRunWarning, "step.sh")
		flow.ExpectNoFile("script-ran.out")
	})
}
