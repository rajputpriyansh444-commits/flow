package command

import (
	"context"
	"testing"

	cmnvalue "github.com/rajputpriyansh444-commits/flow/internal/cmn/value"
	"github.com/rajputpriyansh444-commits/flow/internal/executor/registry"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandExecutorCommandResolutionUsesShellFacts(t *testing.T) {
	ctx := runtime.NewContextForTest(context.Background(), &ir.DAG{Name: "test-dag"}, "run-1", "test.log")
	step := ir.Step{
		Shell:          "direct",
		ExecutorConfig: ir.ExecutorConfig{Type: "command"},
	}
	env := runtime.NewEnv(ctx, step)
	ctx = runtime.WithEnv(ctx, env)
	t.Setenv("FLOW_COMMAND_RESOLUTION_OS", "from-os")

	command := registry.CommandResolution(ctx, step)
	assert.Equal(t, cmnvalue.CommandTargetLocal, command.Target)
	assert.True(t, command.ShellConfigured)
	require.Equal(t, []string{"direct"}, command.Shell)

	got, err := runtime.ResolveString(ctx, "$FLOW_COMMAND_RESOLUTION_OS", cmnvalue.DirectCommandField("command", command))
	require.NoError(t, err)
	assert.Equal(t, "from-os", got)
}
