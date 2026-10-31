//go:build !windows

package procutil

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanLookupStartTimeRejectsPIDsAboveInt32(t *testing.T) {
	t.Parallel()

	if strconv.IntSize < 64 {
		t.Skip("requires an int that can represent a PID above int32")
	}

	one := int64(1)
	require.False(t, canLookupStartTime(int(one<<31)))
}
