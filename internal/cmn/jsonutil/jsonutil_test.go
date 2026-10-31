package jsonutil_test

import (
	"strings"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/jsonutil"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnescaped(t *testing.T) {
	t.Run("PreservesHTMLCharacters", func(t *testing.T) {
		data, err := jsonutil.MarshalUnescaped(map[string]any{
			"text": "a < b & c > d",
		})
		require.NoError(t, err)
		require.Equal(t, `{"text":"a < b & c > d"}`, string(data))
	})
	t.Run("NoTrailingNewline", func(t *testing.T) {
		data, err := jsonutil.MarshalUnescaped([]any{"x", "y"})
		require.NoError(t, err)
		require.Equal(t, `["x","y"]`, string(data))
		require.False(t, strings.HasSuffix(string(data), "\n"))
	})
	t.Run("Error", func(t *testing.T) {
		_, err := jsonutil.MarshalUnescaped(func() {})
		require.Error(t, err)
	})
}
