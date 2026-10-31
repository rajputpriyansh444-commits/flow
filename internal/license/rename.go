package license

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"
)

// MaxServerNameLength is the most characters Flow Console shows of a
// server's name.
const MaxServerNameLength = 100

var (
	// ErrServerNameTooLong is returned when renaming a server to a name
	// longer than MaxServerNameLength.
	ErrServerNameTooLong = errors.New("server names are at most 100 characters")
	// ErrRenameUnsupported is returned when Flow Console cannot rename
	// servers.
	ErrRenameUnsupported = errors.New("renaming servers is not supported by Flow Console yet")
)

// RenameServer sets the name Flow Console shows for this server; an empty
// name shows the name the server reports again. Flow Console keeps the name,
// so it works only for a license that checks in, and fails with
// ErrNoCheckIn otherwise. A *CloudError carries a refusal by Flow Console.
func (m *Manager) RenameServer(ctx context.Context, name string) error {
	name = strings.TrimSpace(name)
	if utf8.RuneCountInString(name) > MaxServerNameLength {
		return ErrServerNameTooLong
	}

	m.transitionMu.Lock()
	defer m.transitionMu.Unlock()

	ad := m.currentActivation()
	claims := m.state.Claims()
	if ad == nil || claims == nil || !m.Source().NeedsHeartbeat() {
		return ErrNoCheckIn
	}
	ctx, cancel := context.WithTimeout(ctx, releaseTimeout)
	defer cancel()
	resp, err := m.client.Rename(ctx, RenameRequest{
		LicenseID:       claims.ID,
		ServerID:        ad.ServerID,
		HeartbeatSecret: ad.HeartbeatSecret,
		Name:            name,
	})
	if err != nil {
		if cloudErr, ok := errors.AsType[*CloudError](err); ok && cloudErr.StatusCode == http.StatusNotFound {
			return ErrRenameUnsupported
		}
		return err
	}
	// Renaming is not a check-in, so the last check-in stays as it was.
	_, err = m.applyToken(ad, resp.Token, ad.CheckedInAt)
	return err
}
