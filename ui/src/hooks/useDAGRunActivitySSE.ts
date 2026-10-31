import { components } from '../api/v1/schema';
import { buildSSEEndpoint, SSEState, useSSE } from './useSSE';

type DAGRunActivity = components['schemas']['DAGRunActivity'];

type DAGRunActivitySSEParams = {
  hours?: number;
  labels?: string;
  workspace?: components['parameters']['Workspace'];
};

/** Streams the activity of a window that ends now. */
export function useDAGRunActivitySSE(
  params: DAGRunActivitySSEParams,
  enabled: boolean = true,
  remoteNode?: string
): SSEState<DAGRunActivity> {
  const endpoint = buildSSEEndpoint('/events/dag-run-activity', params);
  return useSSE<DAGRunActivity>(endpoint, enabled, remoteNode);
}
