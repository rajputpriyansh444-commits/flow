package spec077_xlsx_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/conformance/harness"
	"github.com/stretchr/testify/require"
)

// Spec 077 "Writing" and the append paragraph: rows appended below a header
// row land under the header cells of the same names, a name the header
// lacks adds a column at the right, and a loosely matching name is refused.
func TestXlsxAppendAligned(t *testing.T) {
	t.Parallel()
	flow := harness.NewRunner(t)
	flow.Run("start", "append_aligned.yaml").ExpectExitCode(0)
	var out struct {
		Changes changes          `json:"changes"`
		Headers []string         `json:"headers"`
		Rows    []map[string]any `json:"rows"`
	}
	readJSON(t, flow, "out.json", &out)
	require.Equal(t, changes{Sheet: "Sheet1", Range: "Sheet1!A4:D4", RowsAppended: 1, ColumnsAdded: 1, CellsChanged: 3}, out.Changes)
	require.Equal(t, []string{"when", "what", "who", "note"}, out.Headers)
	require.Len(t, out.Rows, 3)
	require.Equal(t, "2026-10-03", out.Rows[2]["when"], "the field named when lands under the when header, whatever its key order")
	require.Equal(t, "cid", out.Rows[2]["who"])
	require.Nil(t, out.Rows[2]["what"], "a header column no field carries stays empty")
	require.Equal(t, "late", out.Rows[2]["note"], "a field the header lacks adds a column at the right")

	loose := harness.NewRunner(t)
	result := loose.Run("start", "append_loose.yaml")
	result.ExpectNonZeroExitCode()
	result.ExpectStderrContains(`column "amount" not found in header row 1; did you mean "Amount"?`)
	loose.ExpectNoFile("after.txt")
}
