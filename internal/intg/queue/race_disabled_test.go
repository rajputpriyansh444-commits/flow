//go:build !race

package queue_test

func raceEnabled() bool {
	return false
}
