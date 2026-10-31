package profile

import (
	"encoding/hex"
	"fmt"
)

// SecretRef returns the deterministic managed-secret reference for a profile entry.
// The format is "runtime-profiles/<profileName>/key-<hexEncodedKey>".
func SecretRef(profileName, key string) string {
	return fmt.Sprintf("runtime-profiles/%s/key-%s", profileName, hex.EncodeToString([]byte(key)))
}
