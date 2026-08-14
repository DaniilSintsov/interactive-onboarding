import { beforeEach, describe, expect, it, vi } from "vitest";

const mocks = vi.hoisted(() => ({
  effects: [] as Array<() => void | (() => void)>,
  pathname: "/",
  start: vi.fn<(input: { userId: string; testToken?: string }) => Promise<void>>(),
}));

vi.mock("react", () => ({
  createContext: () => ({}),
  use: () => null,
  useCallback: <T extends (...args: never[]) => unknown>(callback: T) => callback,
  useEffect: (effect: () => void | (() => void)) => mocks.effects.push(effect),
  useRef: <T,>(initial: T) => ({ current: initial }),
}));

vi.mock("next/navigation", () => ({
  usePathname: () => mocks.pathname,
}));

vi.mock("@/features/onboarding/model/runtime-path", () => ({
  consumeTestToken: (url: URL) => ({
    token: url.searchParams.get("token") ?? undefined,
    path: url.pathname,
  }),
}));

vi.mock("@/features/onboarding/model/preview-user-store", () => ({
  usePreviewUserStore: (selector: (state: object) => unknown) =>
    selector({
      userId: "experienced-user",
      onboarded: true,
      setOnboarded: vi.fn(),
    }),
}));

vi.mock("@interactive-onboarding/sdk", () => ({
  createOnboarding: () => ({
    start: mocks.start,
    completeCurrentStep: vi.fn(),
    destroy: vi.fn(),
  }),
}));

describe("OnboardingProvider", () => {
  beforeEach(() => {
    mocks.effects.length = 0;
    mocks.pathname = "/";
    mocks.start.mockReset().mockResolvedValue(undefined);
    vi.stubGlobal("window", {
      location: { href: "https://preview.test/?token=test-token" },
      history: { state: null, replaceState: vi.fn() },
    });
  });

  it("продолжает test preview для опытного пользователя после приглашения", async () => {
    const { OnboardingProvider } = await import("./onboarding-provider");
    OnboardingProvider({ children: null, enabled: true });

    mocks.effects[0]?.();
    mocks.effects[1]?.();
    await vi.waitFor(() => {
      expect(mocks.start).toHaveBeenCalledWith({
        userId: "experienced-user",
        testToken: "test-token",
      });
    });

    mocks.pathname = "/add-item/category";
    mocks.effects[2]?.();
    await vi.waitFor(() => {
      expect(mocks.start).toHaveBeenCalledTimes(2);
      expect(mocks.start).toHaveBeenLastCalledWith({
        userId: "experienced-user",
        testToken: undefined,
      });
    });
  });

  it("возобновляет test preview после полной навигации без токена в URL", async () => {
    window.location.href = "https://preview.test/add-item/category";
    const { OnboardingProvider } = await import("./onboarding-provider");
    OnboardingProvider({ children: null, enabled: true });

    mocks.effects[0]?.();
    mocks.effects[1]?.();
    mocks.effects[2]?.();

    await vi.waitFor(() => {
      expect(mocks.start).toHaveBeenCalledWith({
        userId: "experienced-user",
        testToken: undefined,
      });
    });
  });
});
