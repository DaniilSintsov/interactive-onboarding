import { beforeEach, describe, expect, it, vi } from "vitest";

function memoryStorage(): Storage {
  const values = new Map<string, string>();
  return {
    get length() { return values.size; },
    clear: () => values.clear(),
    getItem: (key) => values.get(key) ?? null,
    key: (index) => [...values.keys()][index] ?? null,
    removeItem: (key) => values.delete(key),
    setItem: (key, value) => values.set(key, value),
  };
}

describe("preview user store", () => {
  beforeEach(() => {
    vi.resetModules();
    vi.stubGlobal("localStorage", memoryStorage());
  });

  it("создаёт ID, игнорирует устаревший callback и сбрасывает статус для нового ID", async () => {
    const { usePreviewUserStore } = await import("./preview-user-store");
    const store = usePreviewUserStore.getState();

    store.ensureUser();
    const firstUserId = usePreviewUserStore.getState().userId;
    expect(firstUserId).toMatch(/^[0-9a-f-]{36}$/);
    expect(localStorage.getItem("avito-preview-user-v1")).toContain(firstUserId);

    store.setOnboarded("another-user", true);
    expect(usePreviewUserStore.getState().onboarded).toBe(false);

    store.setOnboarded(firstUserId, true);
    expect(usePreviewUserStore.getState().onboarded).toBe(true);
    expect(localStorage.getItem("avito-preview-user-v1")).toContain('"onboarded":true');

    store.regenerateUser();
    expect(usePreviewUserStore.getState()).toMatchObject({ onboarded: false });
    expect(usePreviewUserStore.getState().userId).not.toBe(firstUserId);
  });
});
