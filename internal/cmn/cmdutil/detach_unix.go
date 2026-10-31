//go:build !windows

package cmdutil

import "syscall"

// DetachAttr returns process attributes that start a command in a session
// of its own, so the end of the caller's terminal or process group does not
// end it. breakawayFromJob applies only on Windows.
func DetachAttr(_ bool) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}
