"use client";

import { createContext, use, useCallback, useEffect, useRef } from "react";
import { usePathname } from "next/navigation";

import { consumeTestToken } from "@/features/onboarding/model/runtime-path";
import { usePreviewUserStore } from "@/features/onboarding/model/preview-user-store";

type OnboardingController = {
  start: (input: { userId: string; testToken?: string }) => Promise<void>;
  completeCurrentStep: () => Promise<void>;
  destroy: () => void;
};

const OnboardingContext = createContext<{
  completeCurrentStep: () => Promise<void>;
} | null>(null);

function takeTestToken() {
  const { token, path } = consumeTestToken(new URL(window.location.href));
  if (path) window.history.replaceState(window.history.state, "", path);
  return token;
}

export function OnboardingProvider({
  children,
  enabled,
}: {
  children: React.ReactNode;
  enabled: boolean;
}) {
  const pathname = usePathname();
  const userId = usePreviewUserStore((state) => state.userId);
  const onboarded = usePreviewUserStore((state) => state.onboarded);
  const setOnboarded = usePreviewUserStore((state) => state.setOnboarded);
  const controllerRef = useRef<OnboardingController | null>(null);
  const readyRef = useRef<Promise<OnboardingController | null> | null>(null);
  const testTokenRef = useRef<string | undefined>(undefined);

  useEffect(() => {
    const testToken = takeTestToken();
    if (testToken) testTokenRef.current = testToken;
  }, [pathname]);

  useEffect(() => {
    if (!enabled || !userId || (onboarded && !testTokenRef.current)) return;

    let cancelled = false;
    let instance: OnboardingController | null = null;

    readyRef.current = import("@interactive-onboarding/sdk").then(
      async ({ createOnboarding }) => {
        if (cancelled) return null;

        const created = createOnboarding({
          projectKey: process.env.NEXT_PUBLIC_PROJECT_KEY ?? "pk_demo_avito",
          runtimeUrl: "/api/runtime",
          onUserStateChange: (user) => setOnboarded(user.userId, user.onboarded),
        });
        instance = created;
        controllerRef.current = created;
        const testToken = testTokenRef.current;
        await created
          .start({ userId, testToken })
          .catch(() => undefined);
        if (testTokenRef.current === testToken) testTokenRef.current = undefined;
        return instance;
      },
    );

    return () => {
      cancelled = true;
      instance?.destroy();
      if (controllerRef.current === instance) controllerRef.current = null;
    };
  }, [enabled, onboarded, setOnboarded, userId]);

  useEffect(() => {
    if (!enabled || !userId || (onboarded && !testTokenRef.current)) return;

    void readyRef.current
      ?.then(async (controller) => {
        if (!controller) return;
        const testToken = testTokenRef.current;
        try {
          await controller.start({ userId, testToken });
        } finally {
          if (testTokenRef.current === testToken) testTokenRef.current = undefined;
        }
      })
      .catch(() => undefined);
  }, [enabled, onboarded, pathname, userId]);

  const completeCurrentStep = useCallback(async () => {
    const controller = controllerRef.current ?? (await readyRef.current);
    await controller?.completeCurrentStep();
  }, []);

  return (
    <OnboardingContext value={{ completeCurrentStep }}>
      {children}
    </OnboardingContext>
  );
}

export function useOnboarding() {
  const context = use(OnboardingContext);
  if (!context) throw new Error("useOnboarding must be used inside OnboardingProvider");
  return context;
}
