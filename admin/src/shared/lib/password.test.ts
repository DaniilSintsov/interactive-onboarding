import { execFile } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import { promisify } from 'node:util';
import { afterEach, beforeAll, describe, expect, it } from 'vitest';
import { matchesAdminPassword } from './password';

const execFileAsync = promisify(execFile);
const generator = fileURLToPath(
  new URL('../../../scripts/generate-admin-credentials.mjs', import.meta.url),
);
let credentials: Record<string, string>;

beforeAll(async () => {
  const { stdout } = await execFileAsync(process.execPath, [generator]);
  credentials = Object.fromEntries(
    stdout
      .trim()
      .split('\n')
      .map((line) => {
        const separator = line.indexOf('=');
        return [line.slice(0, separator), line.slice(separator + 1)];
      }),
  );
});

afterEach(() => {
  delete process.env.ADMIN_PASSWORD_HASH;
});

describe('admin password', () => {
  it('generates deployment-safe credentials accepted by the verifier', async () => {
    expect(credentials.ADMIN_LOGIN_PASSWORD).toHaveLength(32);
    expect(credentials.SESSION_SECRET).toMatch(/^[A-Za-z0-9_-]{43}$/);
    expect(credentials.ADMIN_PASSWORD_HASH).toMatch(
      /^scrypt:v1:32768:8:3:[A-Za-z0-9_-]{22}:[A-Za-z0-9_-]{43}$/,
    );

    process.env.ADMIN_PASSWORD_HASH = credentials.ADMIN_PASSWORD_HASH;
    expect(await matchesAdminPassword(credentials.ADMIN_LOGIN_PASSWORD)).toBe(true);
    expect(await matchesAdminPassword('wrong-password')).toBe(false);
  });

  it('fails closed for invalid input or configuration', async () => {
    expect(await matchesAdminPassword('password')).toBe(false);

    process.env.ADMIN_PASSWORD_HASH = 'not-a-hash';
    expect(await matchesAdminPassword(credentials.ADMIN_LOGIN_PASSWORD)).toBe(false);

    process.env.ADMIN_PASSWORD_HASH = credentials.ADMIN_PASSWORD_HASH.replace(':v1:', ':v2:');
    expect(await matchesAdminPassword(credentials.ADMIN_LOGIN_PASSWORD)).toBe(false);
    expect(await matchesAdminPassword('x'.repeat(257))).toBe(false);
  });
});
