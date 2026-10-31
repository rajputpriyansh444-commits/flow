package worker

import (
	"context"
	"errors"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/rajputpriyansh444-commits/flow/internal/ir"
	"github.com/rajputpriyansh444-commits/flow/internal/proto/convert"
	"github.com/rajputpriyansh444-commits/flow/internal/service/coordinator"
	"github.com/rajputpriyansh444-commits/flow/internal/serviceregistry"
	coordinatorv1 "github.com/rajputpriyansh444-commits/flow/proto/coordinator/v1"
)

type captureCoordinatorClientForTest struct {
	coordinator.Client
	status *ir.DAGRunStatus
	err    error
}

func (c *captureCoordinatorClientForTest) ReportStatus(_ context.Context, req *coordinatorv1.ReportStatusRequest) (*coordinatorv1.ReportStatusResponse, error) {
	c.status, c.err = convert.ProtoToDAGRunStatus(req.Status)
	return &coordinatorv1.ReportStatusResponse{Accepted: true}, nil
}

func (c *captureCoordinatorClientForTest) ReportStatusTo(ctx context.Context, _ serviceregistry.HostInfo, req *coordinatorv1.ReportStatusRequest) (*coordinatorv1.ReportStatusResponse, error) {
	return c.ReportStatus(ctx, req)
}

type captureStatusPusherForTest struct {
	status *ir.DAGRunStatus
}

func (p *captureStatusPusherForTest) Push(_ context.Context, status ir.DAGRunStatus) error {
	copied := status
	p.status = &copied
	return nil
}

// ReportTaskLoadFailureStatusForTest returns the status emitted for a task load failure.
func ReportTaskLoadFailureStatusForTest(ctx context.Context, task *coordinatorv1.Task, root, parent ir.DAGRunRef, loadErr error, profileName string) (*ir.DAGRunStatus, error) {
	client := &captureCoordinatorClientForTest{}
	handler := &remoteTaskHandler{
		workerID:          "worker-test",
		coordinatorClient: client,
	}
	handler.reportTaskLoadFailure(ctx, remoteRun{
		task:        task,
		root:        root,
		parent:      parent,
		profileName: profileName,
	}, loadErr)
	if client.err != nil {
		return nil, client.err
	}
	if client.status == nil {
		return nil, errors.New("load failure status was not reported")
	}
	return client.status, nil
}

// ReportTaskInitFailureStatusForTest returns the status emitted for a task init failure.
func ReportTaskInitFailureStatusForTest(ctx context.Context, task *coordinatorv1.Task, root, parent ir.DAGRunRef, initErr error, profileName string) (*ir.DAGRunStatus, error) {
	pusher := &captureStatusPusherForTest{}
	handler := &remoteTaskHandler{}
	handler.reportTaskInitFailure(ctx, remoteRun{
		task:        task,
		root:        root,
		parent:      parent,
		profileName: profileName,
		handlers:    runHandlers{status: pusher},
	}, initErr)
	if pusher.status == nil {
		return nil, errors.New("init failure status was not reported")
	}
	return pusher.status, nil
}

// LoadRemoteTaskDAGForTest loads the DAG definition for a remote task.
func LoadRemoteTaskDAGForTest(ctx context.Context, cfg *config.Config, task *coordinatorv1.Task) (*ir.DAG, func(), error) {
	if cfg == nil {
		cfg = &config.Config{}
	}
	handler := &remoteTaskHandler{config: cfg}
	loaded, err := handler.loadDAG(ctx, task)
	if err != nil {
		return nil, nil, err
	}
	return loaded.dag, loaded.cleanup, nil
}

// CleanupErrorMessageForTest exposes cleanupErrorMessage for tests.
func CleanupErrorMessageForTest(err error) string {
	return cleanupErrorMessage(err)
}
