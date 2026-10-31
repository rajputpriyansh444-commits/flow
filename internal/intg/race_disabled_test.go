//go:build !race

package intg_test

func raceEnabled() bool {
	return false
}
