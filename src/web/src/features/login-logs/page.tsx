import { useQuery } from '@tanstack/react-query'
import { useNavigate, useSearch } from '@tanstack/react-router'
import { Badge } from '#/components/ui/badge'
import { Card, CardContent } from '#/components/ui/card'
import { adminQueries } from '#/features/admin/queries'
import { DataTable, PageHeader, Pagination, QueryState } from '#/features/admin/shared'

export function LoginLogsPage() {
  const navigate = useNavigate({ from: '/login-logs' })
  const { page } = useSearch({ from: '/_auth/login-logs' })
  const query = useQuery(adminQueries.loginLogs(page))
  return (
    <div className="grid gap-5">
      <PageHeader title="登录日志" subtitle="登录成功、失败和限流审计" />
      <Card>
        <CardContent>
          <QueryState pending={query.isPending} refetching={query.isRefetching} error={query.error} empty={!query.data?.items.length}>
            <DataTable
              headers={['时间', '用户', '事件', '结果', 'IP', '设备', 'User Agent']}
              rows={(query.data?.items ?? []).map((log) => ({
                key: log.id,
                cells: [
                  new Date(log.createdAt).toLocaleString(),
                  log.username,
                  log.event,
                  <Badge variant={log.success ? 'default' : 'destructive'}>{log.success ? '成功' : '失败'}</Badge>,
                  log.ip,
                  log.deviceId || '-',
                  <span className="max-w-96 truncate">{log.userAgent}</span>,
                ],
              }))}
            />
            <Pagination page={page} total={query.data?.total ?? 0} onPageChange={(nextPage) => void navigate({ search: { page: nextPage } })} />
          </QueryState>
        </CardContent>
      </Card>
    </div>
  )
}
