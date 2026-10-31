import { describe, expect, it } from 'vitest';
import { humanizeIdentifier } from '../runtimeConditions';

describe('humanizeIdentifier', () => {
  it('splits acronym boundaries before capitalized words', () => {
    expect(humanizeIdentifier('DAGSnapshotUnavailable')).toBe(
      'DAG Snapshot Unavailable'
    );
    expect(humanizeIdentifier('HTTPServerError')).toBe('HTTP Server Error');
  });
});
