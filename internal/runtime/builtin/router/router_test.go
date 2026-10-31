package router

import (
	"bytes"
	"context"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/runtime"
	"github.com/stretchr/testify/require"
)

func TestRunMasksResolvedSecret(t *testing.T) {
	t.Parallel()

	const secret = "router-secret-value"

	ctx := runtime.NewContext(
		context.Background(),
		&ir.DAG{Name: "router-secret"},
		"run-id",
		"dag.log",
		runtime.WithSecrets([]string{"ROUTE_SECRET=" + secret}),
	)
	exec, err := newRouter(ctx, ir.Step{
		Router: &ir.RouterConfig{
			Value: "${ROUTE_SECRET}",
			Routes: []ir.RouteEntry{
				{Pattern: secret, Targets: []string{secret}},
			},
		},
	})
	require.NoError(t, err)

	var stdout bytes.Buffer
	exec.SetStdout(&stdout)
	require.NoError(t, exec.Run(ctx))

	require.Equal(t, "Router evaluating: *******\n  ******* -> [*******]\n", stdout.String())
}
