import { createFileRoute, redirect } from '@tanstack/react-router'
import { DictionariesPage } from '#/features/dictionaries/page'
import { requirePage } from '#/lib/auth'
export const Route = createFileRoute('/_auth/dictionaries')({ beforeLoad: async ({ context }) => { if (!(await requirePage(context.queryClient, 'dictionaries'))) throw redirect({ to: '/forbidden' }) }, component: DictionariesPage })
