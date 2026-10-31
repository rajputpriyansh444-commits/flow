package spec047_sqlite_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

func TestSQLiteActions(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		file string
		want map[string]int
	}{
		{"positional_params.yaml", map[string]int{"sum": 5}},
		{"import_basic.yaml", map[string]int{"n": 3}},
	} {
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()

			flow := harness.NewRunner(t)
			result := flow.Run("start", tc.file)
			result.ExpectExitCode(0)
			data, err := os.ReadFile(flow.ProjectPath("result.json"))
			require.NoError(t, err)
			var row map[string]int
			require.NoError(t, json.Unmarshal(data, &row))
			require.Equal(t, tc.want, row)
		})
	}
}

func TestSQLiteConfig(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct{ file, command, field string }{
		{"missing_query.yaml", "validate", "with.query"},
		{"missing_import.yaml", "validate", "with.import"},
		{"missing_dsn.yaml", "start", "dsn"},
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
