package spec062_dbt_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

// Validation checks the remote reference without provisioning its runtime.
func TestDBTConfig(t *testing.T) {
	t.Parallel()

	t.Run("versioned action", func(t *testing.T) {
		t.Parallel()

		flow := harness.NewRunner(t)
		flow.Run("validate", "configured.yaml").ExpectExitCode(0)
	})

	t.Run("missing version", func(t *testing.T) {
		t.Parallel()

		flow := harness.NewRunner(t)
		result := flow.Run("validate", "no_version_suffix.yaml")
		result.ExpectNonZeroExitCode()
		result.ExpectStderrContains(`unknown action "dbt"`)
	})
}
