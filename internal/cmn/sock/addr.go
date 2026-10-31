package sock

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/fileutil"
)

// Addr returns the local socket address for an identity and run ID.
func Addr(identity, runID string) string {
	const (
		hashLength          = 6
		maxSocketNameLength = 50
		prefix              = "@flow_"
		suffix              = ".sock"
	)

	hash := fmt.Sprintf("%x", md5.Sum([]byte(identity+runID)))[:hashLength] //nolint:gosec
	safeName := fileutil.SafeName(identity)
	maxNameLen := maxSocketNameLength - len(prefix) - 1 - len(hash) - len(suffix)
	if len(safeName) > maxNameLen {
		safeName = safeName[:maxNameLen]
	}

	return filepath.Join(socketDir(), fmt.Sprintf("%s%s_%s%s", prefix, safeName, hash, suffix))
}

func socketDir() string {
	if runtime.GOOS == "windows" {
		return os.TempDir()
	}
	return "/tmp"
}
