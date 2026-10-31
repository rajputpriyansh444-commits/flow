//go:build !linux && !windows

package resourcelimit

import (
	"context"
	"fmt"
	"runtime"
)

func startNativeGuard(context.Context, Options) (nativeGuard, error) {
	return nil, fmt.Errorf("native resource limits are unsupported on %s", runtime.GOOS)
}
