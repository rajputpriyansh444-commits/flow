package spec043_sftp_transfer_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
)

// Invalid transfer configuration fails without contacting a server.
func TestSFTPConfig(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		file      string
		command   string
		wantError string
	}{
		{"upload.yaml", "validate", ""},
		{"download.yaml", "validate", ""},
		{"direction_mismatch.yaml", "validate", "direction must be"},
		{"missing_source_no_server.yaml", "start", "source path is required"},
		{"missing_destination_no_server.yaml", "start", "destination path is required"},
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
