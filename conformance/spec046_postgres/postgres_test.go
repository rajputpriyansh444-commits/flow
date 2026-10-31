package spec046_postgres_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

func TestPostgresConfig(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		file      string
		command   string
		wantError string
	}{
		{"query.yaml", "validate", ""},
		{"import.yaml", "validate", ""},
		{"missing_query.yaml", "validate", "with.query"},
		{"missing_import.yaml", "validate", "with.import"},
		{"missing_dsn.yaml", "start", "dsn"},
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
