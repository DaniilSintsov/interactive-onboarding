import { NextRequest, NextResponse } from 'next/server';
import { rejectCrossOrigin } from '@/shared/lib/auth';
import { matchesAdminPassword } from '@/shared/lib/password';
import {
  createSessionCookie,
  SESSION_COOKIE,
  SESSION_TTL_SECONDS,
} from '@/shared/lib/session';

export const runtime = 'nodejs';

export async function POST(request: NextRequest) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;

  const body = (await request.json().catch(() => null)) as { password?: unknown } | null;
  if (typeof body?.password !== 'string' || !(await matchesAdminPassword(body.password))) {
    return NextResponse.json({ message: 'Неверный пароль' }, { status: 401 });
  }

  const response = NextResponse.json({ ok: true });
  response.cookies.set(SESSION_COOKIE, await createSessionCookie(), {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    path: '/',
    maxAge: SESSION_TTL_SECONDS,
  });
  return response;
}
