import { NextRequest, NextResponse } from 'next/server';
import { SESSION_COOKIE, verifySessionCookie } from './session';

export async function requireAdmin(request: NextRequest): Promise<NextResponse | null> {
  const session = await verifySessionCookie(request.cookies.get(SESSION_COOKIE)?.value);
  return session ? null : NextResponse.json({ message: 'Требуется вход' }, { status: 401 });
}

export function rejectCrossOrigin(request: NextRequest): NextResponse | null {
  const origin = request.headers.get('origin');
  const host = request.headers.get('host');
  if (!origin || !host) return null;

  try {
    return new URL(origin).host === host
      ? null
      : NextResponse.json({ message: 'Запрос отклонён' }, { status: 403 });
  } catch {
    return NextResponse.json({ message: 'Запрос отклонён' }, { status: 403 });
  }
}
