"use client";

import "@ant-design/v5-patch-for-react-19";

import { useEffect, useState } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ConfigProvider } from "antd";

import { OnboardingProvider } from "@/features/onboarding/ui/onboarding-provider";
import { useDraftStore } from "@/features/listing/model/draft-store";

export function Providers({ children }: { children: React.ReactNode }) {
  const [hydrated, setHydrated] = useState(false);
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: { staleTime: 60_000, retry: 1, refetchOnWindowFocus: false },
        },
      }),
  );

  useEffect(() => {
    let active = true;
    void Promise.resolve(useDraftStore.persist.rehydrate()).then(() => {
      if (active) setHydrated(true);
    });
    return () => {
      active = false;
    };
  }, []);

  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#00aaff",
          colorInfo: "#00aaff",
          colorSuccess: "#159957",
          colorError: "#e83a4f",
          colorText: "#111111",
          colorTextSecondary: "#666666",
          borderRadius: 8,
          borderRadiusLG: 12,
          controlHeight: 44,
          fontFamily: '"Avenir Next", "Helvetica Neue", sans-serif',
        },
        components: {
          Button: { fontWeight: 700, primaryShadow: "none" },
          Card: { boxShadowTertiary: "0 12px 40px rgba(31, 56, 88, .08)" },
          Input: { activeShadow: "0 0 0 3px rgba(0, 111, 251, .13)" },
        },
      }}
    >
      <QueryClientProvider client={queryClient}>
        <OnboardingProvider enabled={hydrated}>{children}</OnboardingProvider>
      </QueryClientProvider>
    </ConfigProvider>
  );
}
