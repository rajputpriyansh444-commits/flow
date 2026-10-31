// Package spec057_redis holds black-box conformance tests for Redis action configuration.
package spec057_redis_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

func TestRedisConfig(t *testing.T) {
	t.Parallel()

	t.Run("accepts configured actions", func(t *testing.T) {
		t.Parallel()

		flow := harness.NewRunner(t)
		flow.Run("validate", "configured.yaml").ExpectExitCode(0)
	})

	for _, tc := range []struct {
		command string
		fixture string
		message string
	}{
		{"validate", "command_mismatch.yaml", "command must be \"GET\" for this action"},
		{"validate", "invalid_mode.yaml", "mode"},
		{"start", "cluster_missing_fields.yaml", "cluster_addrs is required for cluster mode"},
		{"start", "tls_cert_without_key.yaml", "both tls_cert and tls_key must be provided together"},
	} {
		t.Run(tc.fixture, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run(tc.command, tc.fixture)
			result.ExpectNonZeroExitCode()
			result.ExpectStderrContains(tc.message)
		})
	}
}
