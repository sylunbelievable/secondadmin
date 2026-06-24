import { createFileRoute, redirect } from '@tanstack/react-router'
import { LoginForm } from '#/features/auth/login-form'
import { requireAuth } from '#/lib/auth'

export const Route = createFileRoute('/login')({
  beforeLoad: async ({ context }) => {
    if (await requireAuth(context.queryClient)) throw redirect({ to: '/dashboard' })
  },
  component: () => <main className="grid min-h-screen place-items-center bg-muted/40 p-4"><LoginForm /></main>,
})
