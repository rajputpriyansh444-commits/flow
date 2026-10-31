//go:build windows

package sock

import (
	"errors"

	"golang.org/x/sys/windows"
)

// isUnsupportedListenError reports whether AF_UNIX is unavailable.
func isUnsupportedListenError(err error) bool {
	return errors.Is(err, windows.WSAEAFNOSUPPORT) ||
		errors.Is(err, windows.WSAEPROTONOSUPPORT)
}
