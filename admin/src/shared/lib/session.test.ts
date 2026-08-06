import { beforeEach, describe, expect, it } from 'vitest';
import { createSessionCookie, matchesAdminPassword, verifySessionCookie } from './session';

describe('admin session', () => {
  beforeEach(() => {
    process.env.SESSION_SECRET = 'test-session-secret-with-at-least-32-characters';
    process.env.ADMIN_PASSWORD = 'shared-password';
  });

  it('round-trips and expires a signed session', async () => {
    const issuedAt = Date.UTC(2026, 7, 6);
    const cookie = await createSessionCookie(issuedAt);

    expect(await verifySessionCookie(cookie, issuedAt + 1000)).toEqual({
      authenticated: true,
      issuedAt,
    });
    expect(await verifySessionCookie(cookie, issuedAt + 8 * 24 * 60 * 60 * 1000)).toBeNull();
  });

  it('rejects tampering and compares the configured password', async () => {
    const cookie = await createSessionCookie();
    const tampered = `${cookie.slice(0, -1)}${cookie.endsWith('A') ? 'B' : 'A'}`;

    expect(await verifySessionCookie(tampered)).toBeNull();
    expect(await matchesAdminPassword('shared-password')).toBe(true);
    expect(await matchesAdminPassword('wrong-password')).toBe(false);
  });
});
