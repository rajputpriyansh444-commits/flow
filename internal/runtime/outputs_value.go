package runtime

import (
	"maps"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
)

// OutputValuesFromNodes extracts typed DAG/action outputs from runtime nodes.
// Payloads that are not objects are skipped; see NodeData.OutputsValueMap.
func OutputValuesFromNodes(nodes []NodeData) map[string]any {
	outputs := make(map[string]any)
	for _, node := range nodes {
		maps.Copy(outputs, node.OutputsValueMap())
	}
	if len(outputs) == 0 {
		return nil
	}
	return outputs
}

// OutputValuesFromExecNodes extracts typed DAG/action outputs from persisted
// nodes. Payloads that are not objects are skipped; see
// NodeData.OutputsValueMap.
func OutputValuesFromExecNodes(nodes []*ir.Node) map[string]any {
	outputs := make(map[string]any)
	for _, node := range nodes {
		if node == nil || node.OutputsValue == nil {
			continue
		}
		var values map[string]any
		if err := decodeOutputJSON(*node.OutputsValue, &values); err != nil {
			continue
		}
		maps.Copy(outputs, values)
	}
	if len(outputs) == 0 {
		return nil
	}
	return outputs
}
