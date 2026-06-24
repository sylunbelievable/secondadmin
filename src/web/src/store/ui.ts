import { create } from 'zustand'
import { persist } from 'zustand/middleware'

type UIState = {
  theme: 'light' | 'dark'
  sidebarExpanded: boolean
  tabs: Array<{ path: string; title: string }>
  toggleTheme: () => void
  setSidebarExpanded: (expanded: boolean) => void
  openTab: (tab: { path: string; title: string }) => void
  closeTab: (path: string) => void
  resetTabs: () => void
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      theme: 'light',
      sidebarExpanded: true,
      tabs: [{ path: '/dashboard', title: '概览' }],
      toggleTheme: () =>
        set((state) => ({ theme: state.theme === 'light' ? 'dark' : 'light' })),
      setSidebarExpanded: (sidebarExpanded) => set({ sidebarExpanded }),
      openTab: (tab) =>
        set((state) => {
          if (state.tabs.some((item) => item.path === tab.path)) return state
          const dashboard = state.tabs.find((item) => item.path === '/dashboard')
            ?? { path: '/dashboard', title: '概览' }
          return { tabs: [dashboard, ...state.tabs.filter((item) => item.path !== '/dashboard').slice(-8), tab] }
        }),
      closeTab: (path) =>
        set((state) => ({
          tabs: state.tabs.filter((item) => item.path !== path),
        })),
      resetTabs: () => set({ tabs: [{ path: '/dashboard', title: '概览' }] }),
    }),
    { name: 'second-admin-ui' },
  ),
)
