import { queryOptions, type QueryClient } from '@tanstack/react-query'
import { currentMenus, currentUser, responseData, type CurrentMenus } from '#/api'

export const authKeys = {
  me: ['auth', 'me'] as const,
  menus: ['auth', 'menus'] as const,
}

export const authQueries = {
  me: () => queryOptions({
    queryKey: authKeys.me,
    queryFn: () => responseData(currentUser({ throwOnError: true })),
    staleTime: 30_000,
  }),
  menus: () => queryOptions({
    queryKey: authKeys.menus,
    queryFn: () => responseData(currentMenus({ throwOnError: true })),
  }),
}

export async function requireAuth(queryClient: QueryClient) {
  try {
    await queryClient.ensureQueryData(authQueries.me())
  } catch {
    return false
  }
  return true
}

export function can(current: CurrentMenus | undefined, permission: string) {
  return current?.permissions.includes(permission) ?? false
}

export async function requirePage(queryClient: QueryClient, component: string) {
  const current = await queryClient.ensureQueryData(authQueries.menus())
  const stack = [...current.menus]
  while (stack.length) {
    const item = stack.pop()!
    if (item.component === component) return true
    stack.push(...(item.children ?? []))
  }
  return false
}
