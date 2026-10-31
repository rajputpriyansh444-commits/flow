package spec049_data_test

import (
	"os"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

func TestDataConvert(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	result := flow.Run("start", "convert_json_to_yaml.yaml")
	result.ExpectExitCode(0)
	data, err := os.ReadFile(flow.ProjectPath("result.out"))
	require.NoError(t, err)
	require.YAMLEq(t, "name: alice\nage: 30\n", string(data))
}

func TestDataText(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct{ file, want string }{
		{"pick_select_raw.yaml", "alice\n"},
		{"pick_select_null_succeeds.yaml", "\n"},
		{"text.yaml", "hello text"},
	} {
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("start", tc.file)
			result.ExpectExitCode(0)
			flow.ExpectFileContent("result.out", tc.want)
		})
	}
}

func TestDataErrors(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct{ file, command, field string }{
		{"missing_from.yaml", "validate", "from"},
		{"pick_missing_select.yaml", "validate", "select"},
		{"convert_malformed_json_string.yaml", "start", "failed to decode JSON"},
	} {
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run(tc.command, tc.file)
			result.ExpectNonZeroExitCode()
			result.ExpectStderrContains(tc.field)
		})
	}
}
