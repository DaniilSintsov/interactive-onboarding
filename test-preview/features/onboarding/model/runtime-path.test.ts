import { describe, expect, it } from "vitest";

import { backendRuntimePath } from "./runtime-path";

describe("backendRuntimePath", () => {
  it("разрешает только публичный runtime-контракт", () => {
    expect(backendRuntimePath(["scenarios", "resolve"])).toBe(
      "/api/v1/sdk/scenarios/resolve",
    );
    expect(backendRuntimePath(["sessions"])).toBe("/api/v1/sdk/sessions");
    expect(backendRuntimePath(["projects"])).toBeNull();
  });
});
