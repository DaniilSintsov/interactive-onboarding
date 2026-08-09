"use client";

import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

type PreviewUserStore = {
  userId: string;
  onboarded: boolean;
  ensureUser: () => void;
  regenerateUser: () => void;
  setOnboarded: (userId: string, onboarded: boolean) => void;
};

export const usePreviewUserStore = create<PreviewUserStore>()(
  persist(
    (set, get) => ({
      userId: "",
      onboarded: false,
      ensureUser: () => {
        if (!get().userId) set({ userId: crypto.randomUUID() });
      },
      regenerateUser: () => set({ userId: crypto.randomUUID(), onboarded: false }),
      setOnboarded: (userId, onboarded) => {
        if (get().userId === userId) set({ onboarded });
      },
    }),
    {
      name: "avito-preview-user-v1",
      storage: createJSONStorage(() => localStorage),
      skipHydration: true,
      partialize: ({ userId, onboarded }) => ({ userId, onboarded }),
    },
  ),
);
