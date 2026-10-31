package cli_test

import (
	"encoding/json"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

// Nested navigation returns the selected property from the root schema.
func TestSchemaDrillsIntoNestedPath(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	root := flow.Run("schema", "dag")
	root.ExpectExitCode(0)
	var schema struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	require.NoError(t, json.Unmarshal([]byte(root.Stdout()), &schema))
	require.Contains(t, schema.Properties, "steps")

	nested := flow.Run("schema", "dag", "steps")
	nested.ExpectExitCode(0)
	require.JSONEq(t, string(schema.Properties["steps"]), nested.Stdout())
}

// TestSchemaShowsConfigRootFields asserts on "coordinator", a property that
// exists only on the config root schema and not on the DAG schema, so the
// test can't pass against the wrong schema.
func TestSchemaShowsConfigRootFields(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	result := flow.Run("schema", "config")
	result.ExpectExitCode(0)
	require.Contains(t, result.Stdout(), `"coordinator"`)
}

func TestSchemaRejectsUnknownName(t *testing.T) {
	t.Parallel()

	flow := harness.NewRunner(t)
	result := flow.Run("schema", "bogus")
	result.ExpectNonZeroExitCode()
}
