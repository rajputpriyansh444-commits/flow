package runtime_test

import (
	"context"
	"testing"

	cmnvalue "github.com/rajputpriyansh444-commits/flow/internal/cmn/value"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/runctx"
	"github.com/rajputpriyansh444-commits/flow/internal/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Masking rebuilds each message; the images a step attached must come
// through it, or a run with secrets would send the model text alone.
func TestMaskSecretsForProviderKeepsImages(t *testing.T) {
	t.Parallel()

	scope := cmnvalue.NewEnvScope(nil, false).WithEntry("TOKEN", "s3cret", cmnvalue.EnvSourceSecret)
	ctx := runtime.WithDAGContext(context.Background(), runctx.Context{DAG: &ir.DAG{Name: "read"}, EnvScope: scope})

	masked := runtime.MaskSecretsForProvider(ctx, []ir.LLMMessage{{
		Role:    ir.LLMRoleUser,
		Content: "the token is s3cret",
		Images:  []ir.LLMImage{{MediaType: "image/png", Data: []byte{1, 2, 3}}},
	}})

	require.Len(t, masked, 1)
	assert.NotContains(t, masked[0].Content, "s3cret")
	assert.Len(t, masked[0].Images, 1)
}
