//go:build !darwin && !linux && !windows

package wiki

import (
	"os"
	"time"
)

// fileCreationTime falls back to ModTime on unsupported platforms.
func fileCreationTime(_ string, info os.FileInfo) time.Time {
	return info.ModTime()
}
