const encoder = new TextEncoder();
const SESSION_TTL_MS = 7 * 24 * 60 * 60 * 1000;

export const SESSION_COOKIE = 'onboarding_admin_session';

export type SessionData = {
  authenticated: true;
  issuedAt: number;
};

function secret() {
  const value = process.env.SESSION_SECRET;
  return value && value.length >= 32 ? encoder.encode(value) : null;
}

function encodeBase64Url(value: string | Uint8Array): string {
  const bytes = typeof value === 'string' ? encoder.encode(value) : value;
  let binary = '';
  for (const byte of bytes) binary += String.fromCharCode(byte);
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

function decodeBase64Url(value: string) {
  const base64 = value.replace(/-/g, '+').replace(/_/g, '/');
  const binary = atob(base64 + '='.repeat((4 - (base64.length % 4)) % 4));
  return Uint8Array.from(binary, (character) => character.charCodeAt(0));
}

async function key(): Promise<CryptoKey> {
  const value = secret();
  if (!value) throw new Error('SESSION_SECRET must contain at least 32 characters');
  return crypto.subtle.importKey('raw', value, { name: 'HMAC', hash: 'SHA-256' }, false, [
    'sign',
    'verify',
  ]);
}

export async function createSessionCookie(now = Date.now()): Promise<string> {
  const header = encodeBase64Url(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
  const payload = encodeBase64Url(
    JSON.stringify({ authenticated: true, issuedAt: now } satisfies SessionData),
  );
  const body = `${header}.${payload}`;
  const signature = await crypto.subtle.sign('HMAC', await key(), encoder.encode(body));
  return `${body}.${encodeBase64Url(new Uint8Array(signature))}`;
}

export async function verifySessionCookie(
  value: string | undefined,
  now = Date.now(),
): Promise<SessionData | null> {
  if (!value || !secret()) return null;
  const parts = value.split('.');
  if (parts.length !== 3) return null;
  const [header, payload, signature] = parts;

  try {
    const valid = await crypto.subtle.verify(
      'HMAC',
      await key(),
      decodeBase64Url(signature),
      encoder.encode(`${header}.${payload}`),
    );
    if (!valid) return null;

    const data = JSON.parse(new TextDecoder().decode(decodeBase64Url(payload))) as SessionData;
    if (data.authenticated !== true || typeof data.issuedAt !== 'number') return null;
    if (data.issuedAt > now + 5 * 60 * 1000 || now - data.issuedAt > SESSION_TTL_MS) return null;
    return data;
  } catch {
    return null;
  }
}

export async function matchesAdminPassword(candidate: string): Promise<boolean> {
  const configured = process.env.ADMIN_PASSWORD;
  if (!configured) return false;
  const [candidateHash, configuredHash] = await Promise.all([
    crypto.subtle.digest('SHA-256', encoder.encode(candidate)),
    crypto.subtle.digest('SHA-256', encoder.encode(configured)),
  ]);
  const left = new Uint8Array(candidateHash);
  const right = new Uint8Array(configuredHash);
  let difference = 0;
  for (let index = 0; index < left.length; index += 1) difference |= left[index] ^ right[index];
  return difference === 0;
}
