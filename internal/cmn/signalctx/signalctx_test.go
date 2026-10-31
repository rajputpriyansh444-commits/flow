package signalctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOSSignalsDisabled(t *testing.T) {
	base := context.Background()
	require.False(t, OSSignalsDisabled(base))
	require.True(t, OSSignalsDisabled(WithOSSignalsDisabled(base)))
}
