package store

import (
	"crypto/sha256"
	"encoding/hex"
)

// distributedRecordKey is the shared on-disk identifier scheme for the
// distributed control-plane stores (lease, active-run, dispatch). It is a
// hex SHA-256 of input so the key surface stays opaque on disk while still
// being deterministic and collision-free in practice.
func distributedRecordKey(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}
