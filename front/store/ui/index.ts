import { create } from 'zustand'

export interface UIState {
  /** Right activity panel open/closed */
  rightPanelOpen: boolean
  /** Currently active breadcrumb label shown in AppHeader */
  breadcrumb: string

  toggleRightPanel: () => void
  setRightPanelOpen: (open: boolean) => void
  setBreadcrumb: (label: string) => void
}

export const useUIStore = create<UIState>((set) => ({
  rightPanelOpen: true,
  breadcrumb: 'Dashboard',

  toggleRightPanel: () =>
    set((s) => ({ rightPanelOpen: !s.rightPanelOpen })),
  setRightPanelOpen: (open) => set({ rightPanelOpen: open }),
  setBreadcrumb: (label) => set({ breadcrumb: label }),
}))
