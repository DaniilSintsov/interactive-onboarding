import { NextRequest, NextResponse } from 'next/server';
import { rejectCrossOrigin } from '@/shared/lib/auth';
import { SESSION_COOKIE } from '@/shared/lib/session';

export async function POST(request: NextRequest) {
  const forbidden = rejectCrossOrigin(request);
  if (forbidden) return forbidden;

  const response = NextResponse.json({ ok: true });
  response.cookies.set(SESSION_COOKIE, '', {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    path: '/',
    maxAge: 0,
  });
  return response;
}
