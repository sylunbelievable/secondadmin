import { createFileRoute, redirect } from '@tanstack/react-router'
import { z } from 'zod'
import { LoginLogsPage } from '#/features/login-logs/page'
import { requirePage } from '#/lib/auth'

export const Route = createFileRoute('/_auth/login-logs')({
  validateSearch: z.object({ page: z.coerce.number().int().positive().catch(1) }),
  beforeLoad: async ({ context }) => { if (!(await requirePage(context.queryClient, 'login-logs'))) throw redirect({ to: '/forbidden' }) },
  component: LoginLogsPage,
})
