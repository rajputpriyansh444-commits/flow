package spec042_ssh_run_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

func TestSSHConfig(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		file      string
		wantError string
	}{
		{"basic.yaml", ""},
		{"missing_command.yaml", "with.command is required"},
	} {
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("validate", tc.file)
			if tc.wantError == "" {
				result.ExpectExitCode(0)
			} else {
				result.ExpectNonZeroExitCode()
				result.ExpectStderrContains(tc.wantError)
			}
		})
	}
}
