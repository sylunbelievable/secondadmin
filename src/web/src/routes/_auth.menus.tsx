import { createFileRoute, redirect } from '@tanstack/react-router'
import { MenusPage } from '#/features/menus/page'
import { requirePage } from '#/lib/auth'
export const Route = createFileRoute('/_auth/menus')({ beforeLoad: async ({ context }) => { if (!(await requirePage(context.queryClient, 'menus'))) throw redirect({ to: '/forbidden' }) }, component: MenusPage })
