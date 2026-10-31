// Package spec053_archive_test covers archive action wiring through the CLI.
package spec053_archive_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

func TestArchiveRoundTrip(t *testing.T) {
	t.Parallel()
	flow := harness.NewRunner(t)
	flow.WriteFile("source.txt", "hello archive")
	flow.Run("start", "roundtrip.yaml").ExpectExitCode(0)
	flow.ExpectFileContent("extracted/source.txt", "hello archive")
	flow.ExpectNoFile("preview.zip")

	data, err := os.ReadFile(flow.ProjectPath("list.json"))
	require.NoError(t, err)
	var out struct {
		Operation  string `json:"operation"`
		TotalFiles int    `json:"totalFiles"`
		Files      []struct {
			Path string `json:"path"`
		} `json:"files"`
	}
	require.NoError(t, json.Unmarshal(data, &out))
	require.Equal(t, "list", out.Operation)
	require.Equal(t, 1, out.TotalFiles)
	require.Len(t, out.Files, 1)
	require.Equal(t, "source.txt", out.Files[0].Path)
}

func TestArchiveValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ command, file, message string }{
		{"validate", "missing_source.yaml", "source"},
		{"validate", "negative_strip_components.yaml", "minimum"},
		{"start", "missing_destination_for_create.yaml", "destination is required for create"},
	} {
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()
			flow := harness.NewRunner(t)
			result := flow.Run(tc.command, tc.file)
			result.ExpectNonZeroExitCode()
			result.ExpectStderrContains(tc.message)
		})
	}
}
