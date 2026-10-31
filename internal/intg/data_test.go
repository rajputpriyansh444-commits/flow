package intg_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/test"
)

func TestDataConvertAction(t *testing.T) {
	t.Parallel()

	th := test.Setup(t)
	dag := th.DAG(t, `type: graph
steps:
  - id: users
    action: data.convert
    with:
      from: csv
      to: json
      data: |
        name,age
        Alice,30
        Bob,25
    output: USERS_JSON

  - id: first_name
    depends: [users]
    action: data.pick
    with:
      from: json
      select: '.[0].name'
      data: ${USERS_JSON}
      raw: true
    output: FIRST_NAME
`)
	agent := dag.Agent()

	agent.RunSuccess(t)

	dag.AssertLatestStatus(t, ir.Succeeded)
	dag.AssertOutputs(t, map[string]any{
		"FIRST_NAME": "Alice",
	})
}
