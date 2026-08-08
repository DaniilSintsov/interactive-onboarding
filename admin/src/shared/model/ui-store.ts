import { create } from 'zustand';

type UiState = {
  navigationOpen: boolean;
  openNavigation: () => void;
  closeNavigation: () => void;
};

export const useUiStore = create<UiState>((set) => ({
  navigationOpen: false,
  openNavigation: () => set({ navigationOpen: true }),
  closeNavigation: () => set({ navigationOpen: false }),
}));
