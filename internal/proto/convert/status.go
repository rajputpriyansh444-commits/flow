// Package convert provides conversion functions between execution types and proto messages.
package convert

import (
	"encoding/json"
	"fmt"

	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	coordinatorv1 "github.com/rajputpriyansh444-commits/flow/proto/coordinator/v1"
)

// DAGRunStatusToProto converts execution.DAGRunStatus to proto.
func DAGRunStatusToProto(s *ir.DAGRunStatus) (*coordinatorv1.DAGRunStatusProto, error) {
	if s == nil {
		return nil, nil
	}
	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DAGRunStatus: %w", err)
	}
	return &coordinatorv1.DAGRunStatusProto{JsonData: string(data)}, nil
}

// ProtoToDAGRunStatus converts proto to execution.DAGRunStatus.
func ProtoToDAGRunStatus(p *coordinatorv1.DAGRunStatusProto) (*ir.DAGRunStatus, error) {
	if p == nil || p.JsonData == "" {
		return nil, nil
	}
	var s ir.DAGRunStatus
	if err := json.Unmarshal([]byte(p.JsonData), &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DAGRunStatus: %w", err)
	}
	return &s, nil
}
