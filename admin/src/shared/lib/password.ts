import { scrypt, timingSafeEqual } from 'node:crypto';

const SCRYPT_OPTIONS = {
  N: 2 ** 15,
  r: 8,
  p: 3,
  maxmem: 64 * 1024 * 1024,
} as const;
const SALT_BYTES = 16;
const HASH_BYTES = 32;
const MAX_PASSWORD_LENGTH = 256;

function decodeBase64Url(value: string | undefined): Buffer | null {
  if (!value || !/^[A-Za-z0-9_-]+$/.test(value)) return null;
  const decoded = Buffer.from(value, 'base64url');
  return decoded.toString('base64url') === value ? decoded : null;
}

function parsePasswordHash(value: string | undefined): { salt: Buffer; hash: Buffer } | null {
  if (!value) return null;
  const parts = value.split(':');
  if (parts.length !== 7) return null;
  const [algorithm, version, N, r, p, encodedSalt, encodedHash] = parts;
  if (
    algorithm !== 'scrypt' ||
    version !== 'v1' ||
    N !== String(SCRYPT_OPTIONS.N) ||
    r !== String(SCRYPT_OPTIONS.r) ||
    p !== String(SCRYPT_OPTIONS.p)
  ) {
    return null;
  }

  const salt = decodeBase64Url(encodedSalt);
  const hash = decodeBase64Url(encodedHash);
  return salt?.length === SALT_BYTES && hash?.length === HASH_BYTES ? { salt, hash } : null;
}

function deriveKey(password: string, salt: Buffer): Promise<Buffer> {
  return new Promise((resolve, reject) => {
    scrypt(password, salt, HASH_BYTES, SCRYPT_OPTIONS, (error, key) => {
      if (error) reject(error);
      else resolve(key);
    });
  });
}

export async function matchesAdminPassword(candidate: string): Promise<boolean> {
  if (!candidate || candidate.length > MAX_PASSWORD_LENGTH) return false;
  const configured = parsePasswordHash(process.env.ADMIN_PASSWORD_HASH);
  if (!configured) return false;

  try {
    return timingSafeEqual(await deriveKey(candidate, configured.salt), configured.hash);
  } catch {
    return false;
  }
}
