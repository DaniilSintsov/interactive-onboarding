import { describe, expect, it } from "vitest";

import { backendRuntimePath, consumeTestToken } from "./runtime-path";

describe("backendRuntimePath", () => {
  it("разрешает только публичный runtime-контракт", () => {
    expect(backendRuntimePath(["scenarios", "resolve"])).toBe(
      "/api/v1/sdk/scenarios/resolve",
    );
    expect(backendRuntimePath(["sessions"])).toBe("/api/v1/sdk/sessions");
    expect(backendRuntimePath(["events"])).toBe("/api/v1/sdk/events");
    expect(backendRuntimePath(["projects"])).toBeNull();
  });
});

describe("consumeTestToken", () => {
  it("возвращает токен и удаляет его из видимого URL", () => {
    expect(
      consumeTestToken(new URL("https://preview.test/add-item?token=%20secret%20&source=admin#step")),
    ).toEqual({ token: "secret", path: "/add-item?source=admin#step" });
  });

  it("не меняет URL без токена", () => {
    expect(consumeTestToken(new URL("https://preview.test/?source=direct"))).toEqual({
      token: undefined,
    });
  });
});
