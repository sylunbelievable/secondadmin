import { useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Controller, useForm } from 'react-hook-form'
import { z } from 'zod'
import { createRole, deleteRole, errorMessage, setRoleApis, setRoleMenus, updateRole, type Role } from '#/api'
import { Badge } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '#/components/ui/dialog'
import { Field, FieldError, FieldLabel } from '#/components/ui/field'
import { Input } from '#/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '#/components/ui/select'
import { adminKeys, adminQueries } from '#/features/admin/queries'
import { ConfirmButton, DataTable, ErrorAlert, IDsDialog, PageHeader, QueryState, usePermission } from '#/features/admin/shared'
import { authKeys } from '#/lib/auth'

const roleSchema = z.object({ code: z.string().trim().min(1, '请输入编码'), name: z.string().trim().min(1, '请输入名称') })
const editSchema = z.object({ name: z.string().trim().min(1, '请输入名称'), status: z.enum(['0', '1']) })
type RoleForm = z.infer<typeof roleSchema>
type EditForm = z.infer<typeof editSchema>

function EditRoleDialog({ role }: { role: Role }) {
  const [open, setOpen] = useState(false)
  const client = useQueryClient()
  const form = useForm<EditForm>({ resolver: zodResolver(editSchema), values: { name: role.name, status: role.status ? '1' : '0' } })
  const mutation = useMutation({
    mutationFn: (values: EditForm) => updateRole({ path: { id: role.id }, body: { name: values.name, status: Number(values.status) as 0 | 1 }, throwOnError: true }),
    onSuccess: async () => { await Promise.all([client.invalidateQueries({ queryKey: adminKeys.roles }), client.invalidateQueries({ queryKey: authKeys.menus })]); setOpen(false) },
    onError: (error) => form.setError('root', { message: errorMessage(error) }),
  })
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild><Button variant="outline" size="sm">编辑</Button></DialogTrigger>
      <DialogContent>
        <DialogHeader><DialogTitle>编辑角色</DialogTitle><DialogDescription>修改角色名称和状态。</DialogDescription></DialogHeader>
        <form className="grid gap-4" onSubmit={form.handleSubmit((values) => mutation.mutate(values))}>
          <Field data-invalid={!!form.formState.errors.name}><FieldLabel htmlFor={`role-name-${role.id}`}>名称</FieldLabel><Input id={`role-name-${role.id}`} {...form.register('name')} /><FieldError errors={[form.formState.errors.name]} /></Field>
          <Controller name="status" control={form.control} render={({ field }) => <Field><FieldLabel>状态</FieldLabel><Select value={field.value} onValueChange={field.onChange}><SelectTrigger className="w-full"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="1">启用</SelectItem><SelectItem value="0">禁用</SelectItem></SelectContent></Select></Field>} />
          <FieldError errors={[form.formState.errors.root]} />
          <DialogFooter><Button disabled={mutation.isPending}>{mutation.isPending ? '保存中…' : '保存'}</Button></DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

export function RolesPage() {
  const client = useQueryClient()
  const query = useQuery(adminQueries.roles())
  const form = useForm<RoleForm>({ resolver: zodResolver(roleSchema), defaultValues: { code: '', name: '' } })
  const create = useMutation({
    mutationFn: (body: RoleForm) => createRole({ body, throwOnError: true }),
    onSuccess: async () => { form.reset(); await client.invalidateQueries({ queryKey: adminKeys.roles }) },
    onError: (error) => form.setError('root', { message: errorMessage(error) }),
  })
  const remove = useMutation({
    mutationFn: (id: number) => deleteRole({ path: { id }, throwOnError: true }),
    onSuccess: () => Promise.all([client.invalidateQueries({ queryKey: adminKeys.roles }), client.invalidateQueries({ queryKey: authKeys.menus })]),
  })
  const assignAPIs = useMutation({ mutationFn: ({ id, ids }: { id: number; ids: number[] }) => setRoleApis({ path: { id }, body: { ids }, throwOnError: true }) })
  const assignMenus = useMutation({
    mutationFn: ({ id, ids }: { id: number; ids: number[] }) => setRoleMenus({ path: { id }, body: { ids }, throwOnError: true }),
    onSuccess: () => client.invalidateQueries({ queryKey: authKeys.menus }),
  })
  const mayCreate = usePermission('system:role:create')
  const mayUpdate = usePermission('system:role:update')
  const mayDelete = usePermission('system:role:delete')
  const mayAPIs = usePermission('system:role:apis')
  const mayMenus = usePermission('system:role:menus')

  return (
    <div className="grid gap-5">
      <PageHeader title="角色管理" subtitle="API 与菜单授权" />
      {mayCreate && <Card><CardHeader><CardTitle>新建角色</CardTitle></CardHeader><CardContent><form className="grid gap-4 md:grid-cols-3" onSubmit={form.handleSubmit((values) => create.mutate(values))}><Field data-invalid={!!form.formState.errors.code}><FieldLabel htmlFor="role-code">编码</FieldLabel><Input id="role-code" {...form.register('code')} /><FieldError errors={[form.formState.errors.code]} /></Field><Field data-invalid={!!form.formState.errors.name}><FieldLabel htmlFor="role-name">名称</FieldLabel><Input id="role-name" {...form.register('name')} /><FieldError errors={[form.formState.errors.name]} /></Field><div className="flex items-end"><Button className="w-full" disabled={create.isPending}>{create.isPending ? '创建中…' : '创建'}</Button></div><div className="md:col-span-3"><FieldError errors={[form.formState.errors.root]} /></div></form></CardContent></Card>}
      <Card><CardContent><ErrorAlert error={remove.error ?? assignAPIs.error ?? assignMenus.error} /><QueryState pending={query.isPending} refetching={query.isRefetching} error={query.error} empty={!query.data?.length}><DataTable headers={['ID', '编码', '名称', '状态', '操作']} rows={(query.data ?? []).map((role) => ({ key: role.id, cells: [role.id, role.code, role.name, <Badge variant={role.status ? 'default' : 'secondary'}>{role.status ? '启用' : '禁用'}</Badge>, <div className="flex flex-wrap gap-2">{mayUpdate && <EditRoleDialog role={role} />}{mayAPIs && <IDsDialog title={`分配 API：${role.name}`} description="填写角色允许访问的 API ID。" pending={assignAPIs.isPending} trigger={<Button variant="outline" size="sm">API 权限</Button>} onSubmit={async (ids) => { await assignAPIs.mutateAsync({ id: role.id, ids }) }} />}{mayMenus && <IDsDialog title={`分配菜单：${role.name}`} description="填写角色可见的菜单 ID。" pending={assignMenus.isPending} trigger={<Button variant="outline" size="sm">菜单权限</Button>} onSubmit={async (ids) => { await assignMenus.mutateAsync({ id: role.id, ids }) }} />}{mayDelete && role.code !== 'admin' && <ConfirmButton description={`删除角色“${role.name}”，此操作不可撤销。`} pending={remove.isPending} onConfirm={() => remove.mutate(role.id)}>删除</ConfirmButton>}</div>] }))} /></QueryState></CardContent></Card>
    </div>
  )
}
