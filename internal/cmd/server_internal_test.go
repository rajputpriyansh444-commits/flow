package cmd

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/rajputpriyansh444-commits/flow/internal/eventstore"
	"github.com/rajputpriyansh444-commits/flow/internal/license"
	"github.com/rajputpriyansh444-commits/flow/internal/persis"
	"github.com/stretchr/testify/require"
)

// The banner names the web UI URL and the setup step while no admin exists.
func TestStartupBanner(t *testing.T) {
	const url = "http://localhost:8080/flow"
	var out strings.Builder
	startupBanner{version: "v2.1.0", url: url, dagsDir: "/srv/dags", setupRequired: true}.write(&out)
	for _, want := range []string{"Flow v2.1.0 is ready", "Web UI:  " + url + "\n", url + "/setup", "DAGs:    /srv/dags"} {
		require.Contains(t, out.String(), want)
	}

	out.Reset()
	startupBanner{url: url}.write(&out)
	require.NotContains(t, out.String(), "Setup:")
}

func TestStopLocalAgentSessionCleanupHonorsShutdownContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	done := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		result <- stopLocalAgentSessionCleanup(ctx, func() {}, done, nil)
	}()

	select {
	case err := <-result:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(100 * time.Millisecond):
		close(done)
		<-result
		t.Fatal("agent session cleanup did not honor the shutdown context")
	}
}

// The reporter runs whenever a license manager does: so that an
// administrator can choose a level later, and, when configuration fixes off,
// so that Flow Console can be told once that the server no longer reports.
func TestStartCloudReport(t *testing.T) {
	newContext := func(report config.ReportLevel) *Context {
		cfg := &config.Config{Cloud: config.CloudConfig{Report: report}}
		cfg.Paths.DataDir = t.TempDir()
		return &Context{Context: t.Context(), Config: cfg, LicenseManager: license.NewTestManager()}
	}

	for _, level := range []config.ReportLevel{"", config.ReportOff, config.ReportHealth, config.ReportRuns} {
		t.Run("level "+string(level), func(t *testing.T) {
			ctx := newContext(level)
			reporter := startCloudReport(ctx, newCloudReportMonitoring(ctx), newCloudAccess(ctx))
			require.NotNil(t, reporter)
			reporter.Stop()
		})
	}

	t.Run("no license manager", func(t *testing.T) {
		ctx := newContext("")
		ctx.LicenseManager = nil
		monitoring := newCloudReportMonitoring(ctx)
		require.Nil(t, monitoring)
		require.Nil(t, newCloudAccess(ctx))
		require.Nil(t, startCloudReport(ctx, monitoring, nil))
	})
}

// The event store and the run history are prepared at every level: an
// administrator can choose runs at any time, and a reporter below runs
// discards where runs stopped being reported.
func TestCloudReportEvents(t *testing.T) {
	runs := persis.NewDAGRunRepository(nil, nil, persis.DAGRunRepositoryOptions{})
	newContext := func(report config.ReportLevel) *Context {
		cfg := &config.Config{Cloud: config.CloudConfig{Report: report}}
		cfg.Paths.DataDir = t.TempDir()
		return &Context{
			Context:     t.Context(),
			Config:      cfg,
			event:       eventstore.New(nil),
			Persistence: Persistence{DAGRunRepository: runs},
		}
	}

	for _, level := range []config.ReportLevel{config.ReportRuns, "", config.ReportHealth, config.ReportOff} {
		events := cloudReportEvents(newContext(level))
		require.NotNil(t, events.Reader)
		require.Equal(t, runs, events.Statuses)
	}
}
