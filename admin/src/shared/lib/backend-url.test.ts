import { describe, expect, it } from 'vitest';
import { buildBackendUrl } from './backend-url';

describe('buildBackendUrl', () => {
  it('keeps the API boundary, encodes segments and copies filters', () => {
    const result = buildBackendUrl(
      ['projects', 'id with spaces', 'analytics', 'total'],
      '?from=2026-08-01T00%3A00%3A00.000Z',
      'http://localhost:8080',
    );

    expect(result.href).toBe(
      'http://localhost:8080/api/v1/projects/id%20with%20spaces/analytics/total?from=2026-08-01T00%3A00%3A00.000Z',
    );
  });

  it('rejects non-http upstreams', () => {
    expect(() => buildBackendUrl(['projects'], '', 'file:///tmp/api')).toThrow(
      'ADMIN_API_URL must use http or https',
    );
  });
});
