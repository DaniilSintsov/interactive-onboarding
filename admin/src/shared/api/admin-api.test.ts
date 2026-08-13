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
});
