//go:build darwin || windows

package desktop

// awakeReason names the power requests that keep the display on while a
// backend is open, as the operating system lists them.
const awakeReason = "Flow computer step is operating the desktop"
