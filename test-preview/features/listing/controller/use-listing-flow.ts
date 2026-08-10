"use client";

import { useRouter } from "next/navigation";

import { useDraftStore } from "../model/draft-store";

export function useListingFlow() {
  const router = useRouter();
  const resetDraft = useDraftStore((state) => state.resetDraft);

  return {
    go: (path: string) => router.push(path),
    startOver: () => {
      resetDraft();
      router.push("/add-item/category");
    },
  };
}
