"use client";

import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

export type DemoUser = "demo-novice" | "demo-expert";
export type ContactMethod = "calls-and-messages" | "calls" | "messages";

export type ListingDraft = {
  category: string;
  title: string;
  subcategory: string;
  condition: "new" | "used";
  saleType: "personal" | "resale";
  photoName: string;
  description: string;
  hasMultiple: boolean;
  address: string;
  price: string;
  phone: string;
  contactMethod: ContactMethod;
  callDevice: "phone" | "browser";
};

export const initialDraft: ListingDraft = {
  category: "",
  title: "",
  subcategory: "",
  condition: "used",
  saleType: "personal",
  photoName: "",
  description: "",
  hasMultiple: false,
  address: "Москва, ул. Лесная, 7",
  price: "",
  phone: "+7 999 123-45-67",
  contactMethod: "calls-and-messages",
  callDevice: "phone",
};

type DraftStore = ListingDraft & {
  demoUser: DemoUser;
  setDemoUser: (demoUser: DemoUser) => void;
  setField: <K extends keyof ListingDraft>(field: K, value: ListingDraft[K]) => void;
  resetDraft: () => void;
};

export const useDraftStore = create<DraftStore>()(
  persist(
    (set) => ({
      ...initialDraft,
      demoUser: "demo-novice",
      setDemoUser: (demoUser) => set({ demoUser }),
      setField: (field, value) => set({ [field]: value }),
      resetDraft: () => set(initialDraft),
    }),
    {
      name: "avito-listing-draft-v1",
      storage: createJSONStorage(() => sessionStorage),
      skipHydration: true,
      partialize: ({
        category,
        title,
        subcategory,
        condition,
        saleType,
        photoName,
        description,
        hasMultiple,
        address,
        price,
        phone,
        contactMethod,
        callDevice,
        demoUser,
      }) => ({
        category,
        title,
        subcategory,
        condition,
        saleType,
        photoName,
        description,
        hasMultiple,
        address,
        price,
        phone,
        contactMethod,
        callDevice,
        demoUser,
      }),
    },
  ),
);
