import { createFileRoute } from '@tanstack/react-router'
import { StatusPage } from '#/components/status-page'

export const Route = createFileRoute('/_auth/forbidden')({
  component: () => <StatusPage type="forbidden" title="暂无权限" description="当前账号没有访问该页面的权限。" />,
})
