//go:build !windows

package wiki

import (
	"errors"
	"fmt"
	"os"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/fileutil"
	wikimodel "github.com/rajputpriyansh444-commits/flow/internal/wiki"
)

func renameNoReplaceFallback(oldPath, newPath string) error {
	info, err := os.Lstat(oldPath)
	if err != nil {
		return err
	}
	if info.Mode().IsRegular() {
		if err := os.Link(oldPath, newPath); err != nil {
			if errors.Is(err, os.ErrExist) {
				return wikimodel.ErrPageAlreadyExists
			}
			return err
		}
		if err := fileutil.Remove(oldPath); err != nil {
			if rollbackErr := fileutil.Remove(newPath); rollbackErr != nil {
				return fmt.Errorf("failed to remove source after linking destination: %w (destination rollback failed: %v)", err, rollbackErr)
			}
			return err
		}
		return nil
	}

	exists, err := pathExistsNoFollow(newPath)
	if err != nil {
		return err
	}
	if exists {
		return wikimodel.ErrPageAlreadyExists
	}
	return fileutil.Rename(oldPath, newPath)
}
