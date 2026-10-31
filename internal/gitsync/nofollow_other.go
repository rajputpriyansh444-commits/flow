//go:build !unix

package gitsync

import "os"

func openRootFileNoFollow(root *os.Root, name string, flag int, perm os.FileMode) (*os.File, error) {
	return root.OpenFile(name, flag, perm)
}
