import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { deleteSession } from '#/api'
import { Badge } from '#/components/ui/badge'
import { Card, CardContent } from '#/components/ui/card'
import { adminKeys, adminQueries } from '#/features/admin/queries'
import { ConfirmButton, DataTable, ErrorAlert, PageHeader, QueryState, usePermission } from '#/features/admin/shared'

export function SessionsPage() {
  const client = useQueryClient()
  const query = useQuery(adminQueries.sessions())
  const remove = useMutation({
    mutationFn: (id: string) => deleteSession({ path: { id }, throwOnError: true }),
    onSuccess: () => client.invalidateQueries({ queryKey: adminKeys.sessions }),
  })
  const mayDelete = usePermission('system:session:delete')
  return (
    <div className="grid gap-5">
      <PageHeader title="在线设备" subtitle="当前账号的有效会话" />
      <Card>
        <CardContent>
          <ErrorAlert error={remove.error} />
          <QueryState pending={query.isPending} refetching={query.isRefetching} error={query.error} empty={!query.data?.length}>
            <DataTable
              headers={['设备', '方式', '登录时间', '操作']}
              rows={(query.data ?? []).map((session) => ({
                key: session.id,
                cells: [
                  session.deviceId,
                  <Badge variant="outline">{session.authMode}</Badge>,
                  new Date(session.createdAt).toLocaleString(),
                  mayDelete && <ConfirmButton description={`下线设备“${session.deviceId}”，该设备需要重新登录。`} pending={remove.isPending} onConfirm={() => remove.mutate(session.id)}>下线</ConfirmButton>,
                ],
              }))}
            />
          </QueryState>
        </CardContent>
      </Card>
    </div>
  )
}
