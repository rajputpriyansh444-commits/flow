package license

import (
	"slices"

	"github.com/golang-jwt/jwt/v5"
)

// PlanFree is the plan of a license that connects a server to Flow Console
// for monitoring only. It unlocks no paid features, whatever features it
// claims.
const PlanFree = "free"

// LicenseClaims represents the custom JWT claims for a Flow license.
type LicenseClaims struct {
	jwt.RegisteredClaims

	ClaimsVersion int      `json:"cv"`
	Plan          string   `json:"plan"`
	Features      []string `json:"features"`
	ActivationID  string   `json:"activation_id"`
	WarningCode   string   `json:"wc,omitempty"`
	GraceDays     *int     `json:"grace_days,omitempty"`
	// Workspace names the Flow Console workspace that owns the license.
	Workspace string `json:"workspace,omitempty"`
	// ServerName is the name Flow Console shows for the server the token was
	// issued to.
	ServerName string `json:"server_name,omitempty"`
}

// HasFeature returns true if the given feature is included in the license claims.
func (c *LicenseClaims) HasFeature(name string) bool {
	return slices.Contains(c.Features, name)
}

// IsFree reports whether the claims are on the free monitoring plan. Nil
// claims are not.
func (c *LicenseClaims) IsFree() bool {
	return c != nil && c.Plan == PlanFree
}
