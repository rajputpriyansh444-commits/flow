package agent

import (
	"strings"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
)

// createProgressReporter creates the progress reporter
func createProgressReporter(dag *ir.DAG, dagRunID string, params []string) ProgressReporter {
	if dag != nil && dag.Type == ir.TypeAgent {
		display := NewAgentDAGProgressDisplay(dag)
		display.SetDAGRunInfo(dagRunID, strings.Join(params, " "))
		return display
	}
	display := NewSimpleProgressDisplay(dag)
	display.SetDAGRunInfo(dagRunID, strings.Join(params, " "))
	return display
}
