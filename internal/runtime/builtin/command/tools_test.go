package command

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	flowtools "github.com/rajputpriyansh444-commits/flow/internal/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveDeclaredToolCommand(t *testing.T) {
	t.Parallel()

	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	data, err := json.Marshal(&flowtools.Manifest{
		Commands: map[string]flowtools.Command{
			"jq": {
				Name: "jq",
				Path: "/tools/aqua/pkgs/jq",
			},
		},
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(manifestPath, data, 0o600))

	path, ok := resolveDeclaredToolCommand(manifestPath, "jq")

	assert.True(t, ok)
	assert.Equal(t, "/tools/aqua/pkgs/jq", path)
}
