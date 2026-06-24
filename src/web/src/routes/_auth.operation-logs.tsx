import { createFileRoute, redirect } from '@tanstack/react-router'
import { z } from 'zod'
import { OperationLogsPage } from '#/features/operation-logs/page'
import { requirePage } from '#/lib/auth'
export const Route = createFileRoute('/_auth/operation-logs')({
  validateSearch: z.object({ page: z.coerce.number().int().positive().catch(1) }),
  beforeLoad: async ({ context }) => { if (!(await requirePage(context.queryClient, 'operation-logs'))) throw redirect({ to: '/forbidden' }) },
  component: OperationLogsPage,
})
