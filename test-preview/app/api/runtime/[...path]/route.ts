import { NextRequest } from "next/server";

import { backendRuntimePath } from "@/features/onboarding/model/runtime-path";

export const runtime = "nodejs";

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> },
) {
  const backendPath = backendRuntimePath((await params).path);
  if (!backendPath) {
    return Response.json({ message: "Runtime route not found" }, { status: 404 });
  }

  const backendUrl = (process.env.BACKEND_URL ?? "http://localhost:8080").replace(
    /\/$/,
    "",
  );
  const projectKey =
    request.headers.get("x-project-key") ??
    process.env.NEXT_PUBLIC_PROJECT_KEY ??
    "pk_demo_avito";
  const headers = new Headers({
    "Content-Type": request.headers.get("content-type") ?? "application/json",
    "X-Project-Key": projectKey,
  });
  const testToken = request.headers.get("x-scenario-test-token");
  if (testToken) headers.set("X-Scenario-Test-Token", testToken);

  try {
    const upstream = await fetch(`${backendUrl}${backendPath}`, {
      method: "POST",
      headers,
      body: await request.text(),
      cache: "no-store",
    });
    const responseHeaders = new Headers();
    const contentType = upstream.headers.get("content-type");
    if (contentType) responseHeaders.set("Content-Type", contentType);

    return new Response(upstream.status === 204 ? null : upstream.body, {
      status: upstream.status,
      headers: responseHeaders,
    });
  } catch {
    return Response.json(
      { code: "backend_unavailable", message: "Runtime API недоступен" },
      { status: 502 },
    );
  }
}
