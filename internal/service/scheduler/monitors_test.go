package scheduler

import (
	"context"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/rajputpriyansh444-commits/flow/internal/cmn/crypto"
	"github.com/rajputpriyansh444-commits/flow/internal/eventstore"
	notificationmodel "github.com/rajputpriyansh444-commits/flow/internal/notification"
	"github.com/rajputpriyansh444-commits/flow/internal/persis"
	persisfile "github.com/rajputpriyansh444-commits/flow/internal/persis/file"
	notificationservice "github.com/rajputpriyansh444-commits/flow/internal/service/notification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type notificationRoundTripFunc func(*http.Request) (*http.Response, error)

func (fn notificationRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestNotificationServiceUsesConfiguredPublicURL(t *testing.T) {
	t.Parallel()

	type receivedRequest struct {
		body string
		err  error
	}
	received := make(chan receivedRequest, 1)
	httpClient := &http.Client{Transport: notificationRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		received <- receivedRequest{body: string(body), err: err}
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})}

	cfg := &config.Config{
		Paths:  config.PathsConfig{DataDir: t.TempDir()},
		Server: config.Server{PublicURL: "https://flow.example.com/workflows"},
	}
	store := newNotificationTestStore(t, cfg)

	settings, err := notificationmodel.Normalize(&notificationmodel.Settings{
		DAGName: "daily-report",
		Enabled: true,
		Events:  []eventstore.EventType{eventstore.TypeDAGRunFailed},
		Targets: []notificationmodel.Target{{
			ID:      "webhook-1",
			Type:    notificationmodel.ProviderWebhook,
			Enabled: true,
			Webhook: &notificationmodel.WebhookTarget{
				URL:             "https://93.184.216.34/webhook",
				MessageTemplate: "{{run.url}}",
			},
		}},
	}, "tester")
	require.NoError(t, err)
	require.NoError(t, store.Save(context.Background(), settings))

	service := newNotificationService(
		cfg,
		store,
		nil,
		nil,
		notificationservice.WithHTTPClient(httpClient),
	)
	results, err := service.SendTest(
		context.Background(),
		"daily-report",
		"webhook-1",
		eventstore.TypeDAGRunFailed,
	)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.True(t, results[0].Delivered)
	request := <-received
	require.NoError(t, request.err)
	assert.Contains(t, request.body, `"message":"https://flow.example.com/workflows/dag-runs/daily-report/notification-test"`)
}

// Command channels need the server's DAG write and run permissions, and
// their commands never see Flow's own credentials.
func TestNotificationServiceCommandChannels(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shell syntax")
	}

	allowed := map[config.Permission]bool{
		config.PermissionWriteDAGs: true,
		config.PermissionRunDAGs:   true,
	}
	outPath := filepath.Join(t.TempDir(), "env.txt")
	cfg := &config.Config{
		Paths:  config.PathsConfig{DataDir: t.TempDir()},
		Server: config.Server{Permissions: allowed},
		Core: config.Core{BaseEnv: config.NewBaseEnv([]string{
			"PATH=" + os.Getenv("PATH"),
			"FLOW_AUTH_BASIC_PASSWORD=secret",
		})},
	}
	service := newNotificationService(cfg, newNotificationTestStore(t, cfg), nil, nil)
	channel, err := service.SaveChannel(context.Background(), &notificationmodel.Channel{
		Name:    "env",
		Type:    notificationmodel.ProviderCommand,
		Enabled: true,
		Command: &notificationmodel.CommandTarget{Command: "env > " + outPath},
	}, "tester")
	require.NoError(t, err)

	results, err := service.SendChannelTest(context.Background(), channel.ID)
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.True(t, results[0].Delivered, results[0].Error)
	body, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.Contains(t, string(body), "DAG_NAME=")
	assert.NotContains(t, string(body), "FLOW_AUTH_BASIC_PASSWORD")

	for _, permission := range []config.Permission{config.PermissionWriteDAGs, config.PermissionRunDAGs} {
		denied := maps.Clone(allowed)
		denied[permission] = false
		cfg := &config.Config{Server: config.Server{Permissions: denied}}
		assert.False(t, newNotificationService(cfg, nil, nil, nil).CommandChannelsEnabled(), permission)
	}
}

func newNotificationTestStore(t *testing.T, cfg *config.Config) notificationmodel.Store {
	t.Helper()

	key, err := crypto.ResolveKey(cfg.Paths.DataDir)
	require.NoError(t, err)
	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)
	store, err := persisfile.NewNotificationStore(
		persisfile.NewBackend(cfg.Paths).Collection(persis.CollectionNotifications),
		encryptor,
	)
	require.NoError(t, err)
	return store
}
