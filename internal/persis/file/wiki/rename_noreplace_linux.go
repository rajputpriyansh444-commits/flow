//go:build linux

package wiki

import (
	"errors"

	wikimodel "github.com/rajputpriyansh444-commits/flow/internal/wiki"
	"golang.org/x/sys/unix"
)

func renameNoReplace(oldPath, newPath string) error {
	err := unix.Renameat2(unix.AT_FDCWD, oldPath, unix.AT_FDCWD, newPath, unix.RENAME_NOREPLACE)
	if err == nil {
		return nil
	}
	if errors.Is(err, unix.EEXIST) {
		return wikimodel.ErrPageAlreadyExists
	}
	if errors.Is(err, unix.ENOSYS) || errors.Is(err, unix.EINVAL) || errors.Is(err, unix.ENOTSUP) {
		return renameNoReplaceFallback(oldPath, newPath)
	}
	return err
}
