package proc

import (
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/sock"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/stretchr/testify/require"
)

func TestDAGSocketAddr(t *testing.T) {
	t.Parallel()

	ref := ir.NewDAGRunRef("mydag", "run123")
	require.Equal(t, sock.Addr("mydag", "run123"), DAGSocketAddr(ref))
}
