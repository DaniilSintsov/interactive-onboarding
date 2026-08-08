import { beforeEach, describe, expect, it } from 'vitest';
import { createSessionCookie, verifySessionCookie } from './session';

describe('admin session', () => {
  beforeEach(() => {
    process.env.SESSION_SECRET = 'test-session-secret-with-at-least-32-characters';
  });

  it('round-trips and expires a signed session', async () => {
    const issuedAt = Date.UTC(2026, 7, 6);
    const cookie = await createSessionCookie(issuedAt);

    expect(await verifySessionCookie(cookie, issuedAt + 1000)).toEqual({
      authenticated: true,
      issuedAt,
    });
    expect(await verifySessionCookie(cookie, issuedAt + 23 * 60 * 60 * 1000)).not.toBeNull();
    expect(await verifySessionCookie(cookie, issuedAt + 25 * 60 * 60 * 1000)).toBeNull();
  });

  it('rejects tampering and sessions signed with a rotated secret', async () => {
    const cookie = await createSessionCookie();
    const tampered = `${cookie.slice(0, -1)}${cookie.endsWith('A') ? 'B' : 'A'}`;

    expect(await verifySessionCookie(tampered)).toBeNull();
    process.env.SESSION_SECRET = 'rotated-session-secret-with-at-least-32-characters';
    expect(await verifySessionCookie(cookie)).toBeNull();
  });
});
