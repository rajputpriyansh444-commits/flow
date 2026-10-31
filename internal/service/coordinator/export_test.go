package coordinator

import (
	"context"
	"testing"
	"time"

	"github.com/rajputpriyansh444-commits/flow/internal/dispatch"
)

// RuntimeDispatcherConfigForTest exposes the coordinator client config for external tests.
func RuntimeDispatcherConfigForTest(dispatcher dispatch.Dispatcher) (*Config, bool) {
	client, ok := dispatcher.(*clientImpl)
	if !ok {
		return nil, false
	}
	return client.config, true
}

// SetDispatchPollWaitForTest overrides distributed poll wait timing.
func SetDispatchPollWaitForTest(t testing.TB, h *Handler, initial, maxWait time.Duration) {
	t.Helper()
	previousInitial := h.dispatchPollInitialWait
	previousMax := h.dispatchPollMaxWait
	h.dispatchPollInitialWait = initial
	h.dispatchPollMaxWait = maxWait
	t.Cleanup(func() {
		h.dispatchPollInitialWait = previousInitial
		h.dispatchPollMaxWait = previousMax
	})
}

// NotifyDispatchAvailableForTest wakes distributed pollers.
func NotifyDispatchAvailableForTest(h *Handler) {
	h.notifyDispatchAvailable()
}

// ExpireHeldDispatchesForTest fails held dispatches whose wait passed at now.
func ExpireHeldDispatchesForTest(h *Handler, now time.Time) {
	h.expireHeldDispatches(context.Background(), now)
}

// DropOrphanedDispatchesForTest drops the pending dispatches, long pending at
// now, whose runs no longer await them.
func DropOrphanedDispatchesForTest(h *Handler, now time.Time) {
	h.dropOrphanedDispatches(context.Background(), now)
}
