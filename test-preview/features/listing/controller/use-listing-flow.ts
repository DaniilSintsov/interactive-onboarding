"use client";

import { useRouter } from "next/navigation";

import { useDraftStore } from "../model/draft-store";

function withPreview(path: string) {
  if (typeof window === "undefined") return path;
  return new URLSearchParams(window.location.search).get("preview") === "1"
    ? `${path}?preview=1`
    : path;
}

export function useListingFlow() {
  const router = useRouter();
  const resetDraft = useDraftStore((state) => state.resetDraft);

  return {
    go: (path: string) => router.push(withPreview(path)),
    startOver: () => {
      resetDraft();
      router.push(withPreview("/add-item/category"));
    },
  };
}
