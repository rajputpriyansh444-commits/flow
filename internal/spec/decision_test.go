package spec_test

import (
	"context"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestDecisionAction(t *testing.T) {
	t.Parallel()
	_, err := spec.LoadYAML(context.Background(), []byte(`
steps:
  - id: classify
    action: decision.evaluate
    with:
      provider: openrouter
      model: typesafe/jev-1.13
      state: A refund request
      questions:
        refund:
          type: noul
          instructions: Is a refund requested?
  - id: consume
    depends: classify
    action: log.write
    with:
      message: ${steps.classify.outputs.answers}
`))
	require.NoError(t, err)
}
