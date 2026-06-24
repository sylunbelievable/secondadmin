import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'
import { AppShell } from '#/components/layout/app-shell'
import { requireAuth } from '#/lib/auth'

export const Route = createFileRoute('/_auth')({
  beforeLoad: async ({ context }) => {
    if (!(await requireAuth(context.queryClient))) throw redirect({ to: '/login' })
  },
  component: () => <AppShell><Outlet /></AppShell>,
})

