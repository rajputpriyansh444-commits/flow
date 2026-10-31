package file

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	persisfile "github.com/rajputpriyansh444-commits/flow/internal/persis/file"
	"github.com/stretchr/testify/require"
)

func TestNewDependenciesRequiresDAGSettingsStorage(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{}
	_, err := NewDependencies(t.Context(), cfg, persisfile.NewBackend(cfg.Paths))

	require.ErrorContains(t, err, "failed to initialize DAG settings store")
}
