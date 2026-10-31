//go:build !darwin && !(windows && (amd64 || arm64))

package desktop

// OpenElements fails: elements are read through Accessibility on macOS and
// UI Automation on 64-bit Windows.
func OpenElements() (Elements, error) {
	return nil, ErrElementsUnsupported
}
