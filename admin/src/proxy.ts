import { NextRequest, NextResponse } from 'next/server';
import { SESSION_COOKIE, verifySessionCookie } from '@/shared/lib/session';

export async function proxy(request: NextRequest) {
  const authenticated = await verifySessionCookie(request.cookies.get(SESSION_COOKIE)?.value);
  const loginPage = request.nextUrl.pathname === '/login';

  if (loginPage && authenticated) return NextResponse.redirect(new URL('/projects', request.url));
  if (loginPage || authenticated) return NextResponse.next();

  const login = new URL('/login', request.url);
  login.searchParams.set('from', request.nextUrl.pathname);
  return NextResponse.redirect(login);
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico|api/).*)'],
};
