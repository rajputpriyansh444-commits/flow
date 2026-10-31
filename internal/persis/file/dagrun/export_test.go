package dagrun

import "context"

// UpdateLatestAttemptPointerForTest exposes latest-attempt pointer updates to external tests.
func UpdateLatestAttemptPointerForTest(ctx context.Context, statusFile string) error {
	return updateLatestAttemptPointer(ctx, statusFile)
}

// LatestAttemptPointerPathForTest exposes latest-attempt pointer paths to external tests.
func LatestAttemptPointerPathForTest(dagRunsDir string) string {
	return latestAttemptPointerPath(dagRunsDir)
}

// RetryCandidateCacheSizeForTest reports retained retry candidate summaries.
func RetryCandidateCacheSizeForTest(store *Store) int {
	store.retryCandidates.mu.Lock()
	defer store.retryCandidates.mu.Unlock()

	return len(store.retryCandidates.entries)
}
