package license

import "time"

// ActivationData holds the persisted activation state.
type ActivationData struct {
	Token           string `json:"token"`
	HeartbeatSecret string `json:"heartbeat_secret"`
	LicenseKey      string `json:"license_key"`
	ServerID        string `json:"server_id"`
	// Via records how the activation was obtained: ConnectedViaKey or
	// ConnectedViaConsole. Empty means a license key.
	Via string `json:"via,omitempty"`
	// CheckedInAt is the last time Flow Console accepted this activation.
	CheckedInAt time.Time `json:"checked_in_at,omitzero"`
}

// ActivationStore provides persistence for activation data.
type ActivationStore interface {
	Load() (*ActivationData, error)
	Save(data *ActivationData) error
	Remove() error
}
