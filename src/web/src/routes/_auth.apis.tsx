import { createFileRoute, redirect } from '@tanstack/react-router'
import { APIsPage } from '#/features/apis/page'
import { requirePage } from '#/lib/auth'
export const Route = createFileRoute('/_auth/apis')({ beforeLoad: async ({ context }) => { if (!(await requirePage(context.queryClient, 'apis'))) throw redirect({ to: '/forbidden' }) }, component: APIsPage })
