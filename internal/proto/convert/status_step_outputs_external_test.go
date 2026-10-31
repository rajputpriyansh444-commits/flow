package convert_test

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/proto/convert"
	"github.com/stretchr/testify/require"
)

func TestDAGRunStatusProtoRoundTripPreservesStepOutputsValue(t *testing.T) {
	t.Parallel()

	stepOutputsValue := `{"image_tag":"v1.2.3","metadata":"{\"image\":\"api\"}"}`
	original := &ir.DAGRunStatus{
		Name:     "build",
		DAGRunID: "run-1",
		Status:   ir.Running,
		Nodes: []*ir.Node{
			{
				Step:             ir.Step{Name: "publish", ID: "publish"},
				Status:           ir.NodeSucceeded,
				StepOutputsValue: &stepOutputsValue,
			},
		},
	}

	protoStatus, err := convert.DAGRunStatusToProto(original)
	require.NoError(t, err)

	roundTripped, err := convert.ProtoToDAGRunStatus(protoStatus)
	require.NoError(t, err)
	require.Len(t, roundTripped.Nodes, 1)
	require.NotNil(t, roundTripped.Nodes[0].StepOutputsValue)
	require.JSONEq(t, stepOutputsValue, *roundTripped.Nodes[0].StepOutputsValue)
}
