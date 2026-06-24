import { createFileRoute, redirect } from '@tanstack/react-router'
import { SessionsPage } from '#/features/sessions/page'
import { requirePage } from '#/lib/auth'
export const Route = createFileRoute('/_auth/sessions')({ beforeLoad: async ({ context }) => { if (!(await requirePage(context.queryClient, 'sessions'))) throw redirect({ to: '/forbidden' }) }, component: SessionsPage })
