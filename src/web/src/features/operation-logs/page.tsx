import { useQuery } from '@tanstack/react-query'
import { useNavigate, useSearch } from '@tanstack/react-router'
import { Card, CardContent } from '#/components/ui/card'
import { adminQueries } from '#/features/admin/queries'
import { DataTable, PageHeader, Pagination, QueryState } from '#/features/admin/shared'

export function OperationLogsPage() {
  const navigate = useNavigate({ from: '/operation-logs' })
  const { page } = useSearch({ from: '/_auth/operation-logs' })
  const query = useQuery(adminQueries.operationLogs(page))
  return (
    <div className="grid gap-5">
      <PageHeader title="操作日志" subtitle="已认证写请求审计" />
      <Card>
        <CardContent>
          <QueryState pending={query.isPending} refetching={query.isRefetching} error={query.error} empty={!query.data?.items.length}>
            <DataTable
              headers={['时间', '用户', '方法', '路径', '状态', '耗时', '请求 ID']}
              rows={(query.data?.items ?? []).map((log) => ({
                key: log.id,
                cells: [
                  new Date(log.createdAt).toLocaleString(),
                  log.userId,
                  log.method,
                  log.path,
                  log.statusCode,
                  `${log.durationMs}ms`,
                  <code className="text-xs">{log.requestId}</code>,
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
