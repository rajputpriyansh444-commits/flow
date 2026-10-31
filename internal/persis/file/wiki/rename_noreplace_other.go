//go:build !darwin && !linux && !windows

package wiki

func renameNoReplace(oldPath, newPath string) error {
	return renameNoReplaceFallback(oldPath, newPath)
}
