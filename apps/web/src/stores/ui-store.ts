import { create } from 'zustand';

export type AppPage =
  | 'dashboard'
  | 'research'
  | 'markets'
  | 'portfolio'
  | 'strategies'
  | 'agents'
  | 'competitions'
  | 'reports'
  | 'prompt-studio'
  | 'model-router'
  | 'settings';

interface UIState {
  sidebarOpen: boolean;
  currentPage: AppPage;
  toggleSidebar: () => void;
  setPage: (page: AppPage) => void;
}

export const useUI = create<UIState>((set) => ({
  sidebarOpen: true,
  currentPage: 'dashboard',
  toggleSidebar: () => set((s) => ({ sidebarOpen: !s.sidebarOpen })),
  setPage: (page) => set({ currentPage: page }),
}));
