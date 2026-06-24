import { createFileRoute, redirect } from '@tanstack/react-router'
import { z } from 'zod'
import { UsersPage } from '#/features/users/page'
import { requirePage } from '#/lib/auth'
export const Route = createFileRoute('/_auth/users')({
  validateSearch: z.object({ page: z.coerce.number().int().positive().catch(1) }),
  beforeLoad: async ({ context }) => { if (!(await requirePage(context.queryClient, 'users'))) throw redirect({ to: '/forbidden' }) },
  component: UsersPage,
})
