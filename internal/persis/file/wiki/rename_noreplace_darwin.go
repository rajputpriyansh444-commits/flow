//go:build darwin

package wiki

import (
	"errors"

	wikimodel "github.com/rajputpriyansh444-commits/flow/internal/wiki"
	"golang.org/x/sys/unix"
)

func renameNoReplace(oldPath, newPath string) error {
	err := unix.RenameatxNp(unix.AT_FDCWD, oldPath, unix.AT_FDCWD, newPath, unix.RENAME_EXCL)
	if errors.Is(err, unix.EEXIST) {
		return wikimodel.ErrPageAlreadyExists
	}
	if errors.Is(err, unix.ENOTSUP) || errors.Is(err, unix.ENOSYS) || errors.Is(err, unix.EINVAL) {
		return renameNoReplaceFallback(oldPath, newPath)
	}
	return err
}
