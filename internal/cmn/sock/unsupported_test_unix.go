//go:build !windows

package sock

import (
	"net"
	"syscall"
)

// unsupportedListenErrorForTest returns a Unix unsupported-listen error.
func unsupportedListenErrorForTest() error {
	return &net.OpError{
		Op:  "listen",
		Net: "unix",
		Err: syscall.EAFNOSUPPORT,
	}
}
