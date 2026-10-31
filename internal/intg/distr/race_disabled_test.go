//go:build !race

package distr_test

func raceEnabled() bool {
	return false
}
