const runtimePaths: Record<string, string> = {
  "scenarios/resolve": "/api/v1/sdk/scenarios/resolve",
  sessions: "/api/v1/sdk/sessions",
  events: "/api/v1/sdk/events",
};

export function backendRuntimePath(parts: string[]) {
  return runtimePaths[parts.join("/")] ?? null;
}
