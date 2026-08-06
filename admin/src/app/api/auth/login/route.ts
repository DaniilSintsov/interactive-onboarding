import { NextRequest, NextResponse } from 'next/server';
import { rejectCrossOrigin } from '@/shared/lib/auth';
import { createSessionCookie, matchesAdminPassword, SESSION_COOKIE } from '@/shared/lib/session';

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
    maxAge: 7 * 24 * 60 * 60,
  });
  return response;
}
