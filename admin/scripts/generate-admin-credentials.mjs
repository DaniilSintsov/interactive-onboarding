import { randomBytes, scrypt } from 'node:crypto';

const options = { N: 2 ** 15, r: 8, p: 3, maxmem: 64 * 1024 * 1024 };
const password = randomBytes(24).toString('base64url');
const salt = randomBytes(16);
const hash = await new Promise((resolve, reject) => {
  scrypt(password, salt, 32, options, (error, key) => {
    if (error) reject(error);
    else resolve(key);
  });
});

console.log(`ADMIN_LOGIN_PASSWORD=${password}`);
console.log(
  `ADMIN_PASSWORD_HASH=scrypt:v1:${options.N}:${options.r}:${options.p}:${salt.toString('base64url')}:${hash.toString('base64url')}`,
);
console.log(`SESSION_SECRET=${randomBytes(32).toString('base64url')}`);
