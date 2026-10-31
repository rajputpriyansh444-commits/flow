package engine

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/stretchr/testify/require"
)

func TestApplyRunOverridesMergesLabels(t *testing.T) {
	dag := &ir.DAG{
		Labels: ir.NewLabels([]string{"env=prod"}),
	}

	applyRunOverrides(dag, RunOptions{
		Labels: []string{"env=prod", "team=platform"},
	})

	require.Equal(t, []string{"env=prod", "team=platform"}, dag.Labels.Strings())
}
