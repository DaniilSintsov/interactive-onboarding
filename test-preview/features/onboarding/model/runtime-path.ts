const runtimePaths: Record<string, string> = {
  "scenarios/resolve": "/api/v1/sdk/scenarios/resolve",
  sessions: "/api/v1/sdk/sessions",
  events: "/api/v1/sdk/events",
};

export function backendRuntimePath(parts: string[]) {
  return runtimePaths[parts.join("/")] ?? null;
}

export function consumeTestToken(url: URL) {
  const token = url.searchParams.get("token")?.trim() || undefined;
  if (!url.searchParams.has("token")) return { token };

  url.searchParams.delete("token");
  return { token, path: `${url.pathname}${url.search}${url.hash}` };
}
