import { queryOptions } from '@tanstack/react-query'
import {
  dictionaryItems,
  listApis,
  listDictionaries,
  listLoginLogs,
  listMenus,
  listOperationLogs,
  listRoles,
  listSessions,
  listUsers,
  responseData,
} from '#/api'

export const adminKeys = {
  usersAll: ['admin', 'users'] as const,
  users: (page: number) => ['admin', 'users', page] as const,
  roles: ['admin', 'roles'] as const,
  apis: ['admin', 'apis'] as const,
  menus: ['admin', 'menus'] as const,
  dictionaries: ['admin', 'dictionaries'] as const,
  dictionaryItems: (code: string) => ['admin', 'dictionary-items', code] as const,
  loginLogs: (page: number) => ['admin', 'login-logs', page] as const,
  operationLogs: (page: number) => ['admin', 'operation-logs', page] as const,
  sessions: ['admin', 'sessions'] as const,
}

const request = { throwOnError: true } as const

export const adminQueries = {
  users: (page: number) => queryOptions({
    queryKey: adminKeys.users(page),
    queryFn: () => responseData(listUsers({ ...request, query: { page, pageSize: 20 } })),
  }),
  roles: () => queryOptions({
    queryKey: adminKeys.roles,
    queryFn: () => responseData(listRoles(request)),
  }),
  apis: () => queryOptions({
    queryKey: adminKeys.apis,
    queryFn: () => responseData(listApis(request)),
  }),
  menus: () => queryOptions({
    queryKey: adminKeys.menus,
    queryFn: () => responseData(listMenus(request)),
  }),
  dictionaries: () => queryOptions({
    queryKey: adminKeys.dictionaries,
    queryFn: () => responseData(listDictionaries(request)),
  }),
  dictionaryItems: (code: string) => queryOptions({
    queryKey: adminKeys.dictionaryItems(code),
    queryFn: () => responseData(dictionaryItems({ ...request, path: { code } })),
    enabled: !!code,
  }),
  operationLogs: (page: number) => queryOptions({
    queryKey: adminKeys.operationLogs(page),
    queryFn: () => responseData(listOperationLogs({ ...request, query: { page, pageSize: 20 } })),
  }),
  loginLogs: (page: number) => queryOptions({
    queryKey: adminKeys.loginLogs(page),
    queryFn: () => responseData(listLoginLogs({ ...request, query: { page, pageSize: 20 } })),
  }),
  sessions: () => queryOptions({
    queryKey: adminKeys.sessions,
    queryFn: () => responseData(listSessions(request)),
  }),
}
