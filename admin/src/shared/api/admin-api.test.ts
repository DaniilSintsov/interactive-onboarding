import { afterEach, describe, expect, it, vi } from 'vitest';
import { adminApi } from './admin-api';

describe('adminApi', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('adds page query when loading project elements for selected page', async () => {
    const fetchMock = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValue(new Response(JSON.stringify([]), { status: 200 }));

    await adminApi.listElements('project-1', '/add-item/details');

    expect(fetchMock).toHaveBeenCalledWith('/api/backend/projects/project-1/elements?page=%2Fadd-item%2Fdetails', {
      headers: {},
    });
  });

  it('returns page paths from project pages response', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ items: [{ page: '/add-item/details' }] }), { status: 200 }),
    );

    await expect(adminApi.listPages('project-1')).resolves.toEqual(['/add-item/details']);
  });

  it('sends page in createElement body', async () => {
    const fetchMock = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValue(new Response(JSON.stringify({ id: 'element-1' }), { status: 200 }));

    await adminApi.createElement('project-1', {
      key: 'publish-ad-button',
      label: 'Кнопка публикации объявления',
      page: '/add-item/details',
    });

    expect(fetchMock).toHaveBeenCalledWith('/api/backend/projects/project-1/elements', {
      method: 'POST',
      body: JSON.stringify({
        key: 'publish-ad-button',
        label: 'Кнопка публикации объявления',
        page: '/add-item/details',
      }),
      headers: { 'content-type': 'application/json' },
    });
  });

  it('sends delete request for project removal', async () => {
    const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(new Response(null, { status: 204 }));

    await adminApi.deleteProject('project-1');

    expect(fetchMock).toHaveBeenCalledWith('/api/backend/projects/project-1', {
      method: 'DELETE',
      headers: {},
    });
  });

  it('uses total analytics endpoint with period query for scenario summary', async () => {
    const fetchMock = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValue(new Response(JSON.stringify({ scenario_id: 'scenario-1' }), { status: 200 }));

    await adminApi.getScenarioAnalyticsTotal(
      'scenario-1',
      '?from=2026-08-01T00%3A00%3A00.000Z&to=2026-08-08T00%3A00%3A00.000Z',
    );

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/backend/scenarios/scenario-1/analytics/total?from=2026-08-01T00%3A00%3A00.000Z&to=2026-08-08T00%3A00%3A00.000Z',
      { headers: {} },
    );
  });
});
