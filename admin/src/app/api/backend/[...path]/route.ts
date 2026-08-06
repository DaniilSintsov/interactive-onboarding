import { NextRequest, NextResponse } from 'next/server';
import { rejectCrossOrigin, requireAdmin } from '@/shared/lib/auth';
import { buildBackendUrl } from '@/shared/lib/backend-url';

const MUTATING_METHODS = new Set(['POST', 'PUT', 'PATCH', 'DELETE']);

async function proxy(
  request: NextRequest,
  context: { params: Promise<{ path: string[] }> },
): Promise<NextResponse> {
  const unauthorized = await requireAdmin(request);
  if (unauthorized) return unauthorized;
  if (MUTATING_METHODS.has(request.method)) {
    const forbidden = rejectCrossOrigin(request);
    if (forbidden) return forbidden;
  }

  const password = process.env.ADMIN_PASSWORD;
  if (!password) {
    return NextResponse.json({ message: 'ADMIN_PASSWORD не настроен' }, { status: 500 });
  }

  try {
    const { path } = await context.params;
    const upstream = await fetch(buildBackendUrl(path, request.nextUrl.search), {
      method: request.method,
      headers: {
        accept: request.headers.get('accept') || 'application/json',
        authorization: `Bearer ${password}`,
        ...(request.headers.get('content-type')
          ? { 'content-type': request.headers.get('content-type') as string }
          : {}),
      },
      body: request.method === 'GET' || request.method === 'HEAD' ? undefined : await request.arrayBuffer(),
      cache: 'no-store',
    });

    const headers = new Headers();
    for (const name of ['content-type', 'content-disposition']) {
      const value = upstream.headers.get(name);
      if (value) headers.set(name, value);
    }
    return new NextResponse(upstream.body, { status: upstream.status, headers });
  } catch (error) {
    console.error('Admin BFF error', error);
    return NextResponse.json({ message: 'Go API недоступен' }, { status: 502 });
  }
}

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const PATCH = proxy;
export const DELETE = proxy;
