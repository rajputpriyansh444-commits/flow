package spec038_kubernetes_run_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

// Configuration checks run without cluster access.
func TestKubernetesConfig(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		file      string
		command   string
		wantError string
	}{
		{"create_basic.yaml", "validate", ""},
		{"missing_target_bare.yaml", "start", "image is required"},
		{"missing_target_with_other.yaml", "validate", "image"},
		{"invalid_cleanup_policy.yaml", "start", "cleanup_policy must be either delete or keep"},
		{"negative_active_deadline.yaml", "validate", "active_deadline"},
	} {
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run(tc.command, tc.file)
			if tc.wantError == "" {
				result.ExpectExitCode(0)
			} else {
				result.ExpectNonZeroExitCode()
				result.ExpectStderrContains(tc.wantError)
			}
		})
	}
}
