import { useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useNavigate, useSearch } from '@tanstack/react-router'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { createUser, errorMessage, setUserRoles, updateUser, type User } from '#/api'
import { Badge } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '#/components/ui/dialog'
import { Field, FieldError, FieldLabel } from '#/components/ui/field'
import { Input } from '#/components/ui/input'
import { adminKeys, adminQueries } from '#/features/admin/queries'
import { ConfirmButton, DataTable, ErrorAlert, IDsDialog, PageHeader, Pagination, QueryState, usePermission } from '#/features/admin/shared'
import { authKeys } from '#/lib/auth'

const createSchema = z.object({
  username: z.string().trim().min(1, '请输入用户名'),
  nickname: z.string().trim(),
  password: z.string().min(8, '密码至少 8 位'),
})
const editSchema = z.object({
  nickname: z.string().trim(),
  password: z.string().refine((value) => value === '' || value.length >= 8, '密码至少 8 位'),
})
type CreateForm = z.infer<typeof createSchema>
type EditForm = z.infer<typeof editSchema>

function EditUserDialog({ user }: { user: User }) {
  const [open, setOpen] = useState(false)
  const client = useQueryClient()
  const form = useForm<EditForm>({
    resolver: zodResolver(editSchema),
    values: { nickname: user.nickname, password: '' },
  })
  const mutation = useMutation({
    mutationFn: (values: EditForm) => updateUser({
      path: { id: user.id },
      body: { nickname: values.nickname, ...(values.password ? { password: values.password } : {}) },
            throwOnError: true,
    }),
    onSuccess: async () => {
      await client.invalidateQueries({ queryKey: adminKeys.usersAll })
      setOpen(false)
    },
    onError: (error) => form.setError('root', { message: errorMessage(error) }),
  })
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild><Button variant="outline" size="sm">编辑</Button></DialogTrigger>
      <DialogContent>
        <DialogHeader><DialogTitle>编辑用户</DialogTitle><DialogDescription>更新昵称或重置密码。</DialogDescription></DialogHeader>
        <form className="grid gap-4" onSubmit={form.handleSubmit((values) => mutation.mutate(values))}>
          <Field><FieldLabel htmlFor={`nickname-${user.id}`}>昵称</FieldLabel><Input id={`nickname-${user.id}`} {...form.register('nickname')} /></Field>
          <Field data-invalid={!!form.formState.errors.password}><FieldLabel htmlFor={`password-${user.id}`}>新密码</FieldLabel><Input id={`password-${user.id}`} type="password" placeholder="留空则不修改" aria-invalid={!!form.formState.errors.password} {...form.register('password')} /><FieldError errors={[form.formState.errors.password]} /></Field>
          <FieldError errors={[form.formState.errors.root]} />
          <DialogFooter><Button disabled={mutation.isPending}>{mutation.isPending ? '保存中…' : '保存'}</Button></DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

export function UsersPage() {
  const client = useQueryClient()
  const navigate = useNavigate({ from: '/users' })
  const { page } = useSearch({ from: '/_auth/users' })
  const query = useQuery(adminQueries.users(page))
  const form = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: { username: '', nickname: '', password: '' },
  })
  const create = useMutation({
    mutationFn: (body: CreateForm) => createUser({ body, throwOnError: true }),
    onSuccess: async () => {
      form.reset()
      await client.invalidateQueries({ queryKey: adminKeys.usersAll })
    },
    onError: (error) => form.setError('root', { message: errorMessage(error) }),
  })
  const update = useMutation({
    mutationFn: ({ id, status }: { id: string; status: 0 | 1 }) => updateUser({
      path: { id }, body: { status }, throwOnError: true,
    }),
    onSuccess: () => client.invalidateQueries({ queryKey: adminKeys.usersAll }),
  })
  const assign = useMutation({
    mutationFn: ({ id, ids }: { id: string; ids: string[] }) => setUserRoles({
      path: { id }, body: { ids }, throwOnError: true,
    }),
    onSuccess: () => client.invalidateQueries({ queryKey: authKeys.menus }),
  })
  const mayCreate = usePermission('system:user:create')
  const mayUpdate = usePermission('system:user:update')
  const mayAssign = usePermission('system:user:roles')

  return (
    <div className="grid gap-5">
      <PageHeader title="用户管理" subtitle="账号、状态与角色" />
      {mayCreate && (
        <Card>
          <CardHeader><CardTitle>新建用户</CardTitle></CardHeader>
          <CardContent>
            <form className="grid gap-4 md:grid-cols-4" onSubmit={form.handleSubmit((values) => create.mutate(values))}>
              <Field data-invalid={!!form.formState.errors.username}><FieldLabel htmlFor="new-username">用户名</FieldLabel><Input id="new-username" aria-invalid={!!form.formState.errors.username} {...form.register('username')} /><FieldError errors={[form.formState.errors.username]} /></Field>
              <Field><FieldLabel htmlFor="new-nickname">昵称</FieldLabel><Input id="new-nickname" {...form.register('nickname')} /></Field>
              <Field data-invalid={!!form.formState.errors.password}><FieldLabel htmlFor="new-password">初始密码</FieldLabel><Input id="new-password" type="password" aria-invalid={!!form.formState.errors.password} {...form.register('password')} /><FieldError errors={[form.formState.errors.password]} /></Field>
              <div className="flex items-end"><Button className="w-full" disabled={create.isPending}>{create.isPending ? '创建中…' : '创建'}</Button></div>
              <div className="md:col-span-4"><FieldError errors={[form.formState.errors.root]} /></div>
            </form>
          </CardContent>
        </Card>
      )}
      <Card>
        <CardContent>
          <ErrorAlert error={update.error ?? assign.error} />
          <QueryState pending={query.isPending} refetching={query.isRefetching} error={query.error} empty={!query.data?.items.length}>
            <DataTable
              headers={['ID', '用户名', '昵称', '状态', '操作']}
              rows={(query.data?.items ?? []).map((user) => ({
                key: user.id,
                cells: [
                  user.id,
                  user.username,
                  user.nickname,
                  <Badge variant={user.status ? 'default' : 'secondary'}>{user.status ? '启用' : '禁用'}</Badge>,
                  <div className="flex flex-wrap gap-2">
                    {mayUpdate && <EditUserDialog user={user} />}
                    {mayUpdate && (user.status
                      ? <ConfirmButton description={`禁用用户“${user.username}”，其后续请求将被拒绝。`} pending={update.isPending} onConfirm={() => update.mutate({ id: user.id, status: 0 })}>禁用</ConfirmButton>
                      : <Button variant="outline" size="sm" disabled={update.isPending} onClick={() => update.mutate({ id: user.id, status: 1 })}>启用</Button>)}
                    {mayAssign && <IDsDialog title={`分配角色：${user.username}`} description="填写该用户拥有的角色 ID。" pending={assign.isPending} trigger={<Button variant="outline" size="sm">分配角色</Button>} onSubmit={async (ids) => { await assign.mutateAsync({ id: user.id, ids }) }} />}
                  </div>,
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
