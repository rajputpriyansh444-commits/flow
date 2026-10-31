//go:build unix

package tools

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// A subprocess bounds the regression: opening a FIFO before stat would hang.
func TestToolsetHashFIFO(t *testing.T) {
	if path := os.Getenv("FLOW_TEST_HASH_FIFO"); path != "" {
		_, _, err := ToolsetHash(&ir.ToolConfig{Registry: &ir.ToolRegistry{Type: "local", Path: path}}, "linux/amd64")
		require.ErrorContains(t, err, "not a regular file")
		return
	}
	path := filepath.Join(t.TempDir(), "registry.fifo")
	require.NoError(t, unix.Mkfifo(path, 0o600))
	binary, err := os.Executable()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, "-test.run=^TestToolsetHashFIFO$")
	cmd.Env = append(os.Environ(), "FLOW_TEST_HASH_FIFO="+path)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", output)
}

func TestToolsetHashDevice(t *testing.T) {
	for _, path := range []string{"/dev/null", "/dev/zero"} {
		t.Run(path, func(t *testing.T) {
			_, _, err := ToolsetHash(&ir.ToolConfig{Registry: &ir.ToolRegistry{Type: "local", Path: path}}, "linux/amd64")
			require.ErrorContains(t, err, "not a regular file")
		})
	}
}

// Symlinks to ordinary registry files remain supported.
func TestToolsetHashSymlink(t *testing.T) {
	path := filepath.Join(t.TempDir(), "registry.yaml")
	require.NoError(t, os.WriteFile(path, []byte("packages: []"), 0o600))
	link := path + ".link"
	require.NoError(t, os.Symlink(path, link))
	_, _, err := ToolsetHash(&ir.ToolConfig{Registry: &ir.ToolRegistry{Type: "local", Path: link}}, "linux/amd64")
	require.NoError(t, err)
}
