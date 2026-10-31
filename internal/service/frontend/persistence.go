package frontend

import (
	"github.com/rajputpriyansh444-commits/flow/internal/audit"
	authmodel "github.com/rajputpriyansh444-commits/flow/internal/auth"
	"github.com/rajputpriyansh444-commits/flow/internal/dagpin"
	"github.com/rajputpriyansh444-commits/flow/internal/dagsettings"
	"github.com/rajputpriyansh444-commits/flow/internal/eventstore"
	"github.com/rajputpriyansh444-commits/flow/internal/incident"
	"github.com/rajputpriyansh444-commits/flow/internal/notification"
	"github.com/rajputpriyansh444-commits/flow/internal/profile"
	"github.com/rajputpriyansh444-commits/flow/internal/publicwebhook"
	"github.com/rajputpriyansh444-commits/flow/internal/remotenode"
	"github.com/rajputpriyansh444-commits/flow/internal/secret"
	authservice "github.com/rajputpriyansh444-commits/flow/internal/service/auth"
	"github.com/rajputpriyansh444-commits/flow/internal/service/chatbridge"
	"github.com/rajputpriyansh444-commits/flow/internal/upgrade"
	"github.com/rajputpriyansh444-commits/flow/internal/view"
	"github.com/rajputpriyansh444-commits/flow/internal/wiki"
	"github.com/rajputpriyansh444-commits/flow/internal/workspace"
)

// Stores contains the persistence dependencies used by the frontend server.
type Stores struct {
	WorkspaceBaseConfig  dagsettings.BaseConfigProvider
	BaseConfig           dagsettings.BaseConfigStore
	AuthService          *authservice.Service
	UserStore            authmodel.UserStore
	AuthSetupRequired    bool
	RemoteNode           remotenode.Store
	Secret               secret.Store
	Profile              profile.Store
	DAGSettings          dagsettings.Store
	DAGPins              dagpin.Store
	Wiki                 wiki.PageStore
	Notification         notification.Store
	NotificationState    chatbridge.StateStore
	NewNotificationLease func() chatbridge.Lease
	Incident             incident.Store
	IncidentState        chatbridge.StateStore
	NewIncidentLease     func() chatbridge.Lease
	Workspace            workspace.Store
	Upgrade              upgrade.CacheStore
	Audit                audit.Store
	Event                *eventstore.Service
	View                 view.Store
	PublicWebhooks       publicwebhook.Store
}
