"use client";

import { createContext, use, useCallback, useEffect, useRef } from "react";
import { usePathname } from "next/navigation";

import { useDraftStore } from "@/features/listing/model/draft-store";

type OnboardingController = {
  start: (input: { userId: string; preview?: boolean }) => Promise<void>;
  completeCurrentStep: () => Promise<void>;
  destroy: () => void;
};

const OnboardingContext = createContext<{
  completeCurrentStep: () => Promise<void>;
} | null>(null);

function isPreview() {
  return new URLSearchParams(window.location.search).get("preview") === "1";
}

export function OnboardingProvider({
  children,
  enabled,
}: {
  children: React.ReactNode;
  enabled: boolean;
}) {
  const pathname = usePathname();
  const demoUser = useDraftStore((state) => state.demoUser);
  const controllerRef = useRef<OnboardingController | null>(null);
  const readyRef = useRef<Promise<OnboardingController | null> | null>(null);

  useEffect(() => {
    if (!enabled) return;

    let cancelled = false;
    let instance: OnboardingController | null = null;

    readyRef.current = import("@interactive-onboarding/sdk").then(
      async ({ createOnboarding }) => {
        if (cancelled) return null;

        const created = createOnboarding({
          projectKey: process.env.NEXT_PUBLIC_PROJECT_KEY ?? "pk_demo_avito",
          runtimeUrl: "/api/runtime",
        });
        instance = created;
        controllerRef.current = created;
        await created
          .start({ userId: demoUser, preview: isPreview() })
          .catch(() => undefined);
        return instance;
      },
    );

    return () => {
      cancelled = true;
      instance?.destroy();
      if (controllerRef.current === instance) controllerRef.current = null;
    };
  }, [demoUser, enabled]);

  useEffect(() => {
    if (!enabled) return;

    void readyRef.current
      ?.then((controller) =>
        controller?.start({ userId: demoUser, preview: isPreview() }),
      )
      .catch(() => undefined);
  }, [demoUser, enabled, pathname]);

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
