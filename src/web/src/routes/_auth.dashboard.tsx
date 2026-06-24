import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card'

export const Route = createFileRoute('/_auth/dashboard')({
  component: () => (
    <div className="grid gap-5">
      <div>
        <h1 className="text-2xl font-semibold">概览</h1>
        <p className="text-muted-foreground">系统状态与快捷入口</p>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Second Admin</CardTitle>
          <CardDescription>管理后台已连接</CardDescription>
        </CardHeader>
        <CardContent>使用左侧菜单进入系统管理功能。</CardContent>
      </Card>
    </div>
  ),
})
