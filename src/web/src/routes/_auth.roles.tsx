import { createFileRoute, redirect } from '@tanstack/react-router'
import { RolesPage } from '#/features/roles/page'
import { requirePage } from '#/lib/auth'
export const Route = createFileRoute('/_auth/roles')({ beforeLoad: async ({ context }) => { if (!(await requirePage(context.queryClient, 'roles'))) throw redirect({ to: '/forbidden' }) }, component: RolesPage })
