// Package proc owns process liveness and control addressing.
package proc

import (
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/sock"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
)

// DAGSocketAddr returns the control socket address for a DAG run.
func DAGSocketAddr(dagRun ir.DAGRunRef) string {
	return sock.Addr(dagRun.Name, dagRun.ID)
}
