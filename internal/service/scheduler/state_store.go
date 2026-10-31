package scheduler

import (
	"context"

	"github.com/rajputpriyansh444-commits/flow/internal/schedulerstate"
)

type noopStateStore struct{}

var _ schedulerstate.Store = noopStateStore{}

func (noopStateStore) Load(_ context.Context) (*schedulerstate.State, error) {
	return &schedulerstate.State{DAGs: make(map[string]schedulerstate.DAGWatermark)}, nil
}

func (noopStateStore) Save(_ context.Context, _ *schedulerstate.State) error {
	return nil
}
