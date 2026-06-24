import { useState, type ReactNode } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Controller, useForm } from 'react-hook-form'
import { z } from 'zod'
import { createMenu, deleteMenu, errorMessage, updateMenu, type Menu } from '#/api'
import { Badge } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '#/components/ui/dialog'
import { Field, FieldError, FieldLabel } from '#/components/ui/field'
import { Input } from '#/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '#/components/ui/select'
import { adminKeys, adminQueries } from '#/features/admin/queries'
import { ConfirmButton, DataTable, ErrorAlert, PageHeader, QueryState, usePermission } from '#/features/admin/shared'
import { authKeys } from '#/lib/auth'

const schema = z.object({
  parentId: z.string().trim().regex(/^\d+$/, '请输入数字 ID'),
  type: z.enum(['directory', 'menu', 'button']),
  name: z.string().trim().min(1, '请输入名称'),
  path: z.string(),
  component: z.string(),
  icon: z.string(),
  permission: z.string(),
  sort: z.number().int(),
  visible: z.enum(['true', 'false']),
  status: z.enum(['0', '1']),
})
type Form = z.infer<typeof schema>

const defaults: Form = { parentId: '0', type: 'menu', name: '', path: '', component: '', icon: '', permission: '', sort: 0, visible: 'true', status: '1' }

function menuValues(menu: Menu): Form {
  return {
    parentId: menu.parentId,
    type: menu.type,
    name: menu.name,
    path: menu.path ?? '',
    component: menu.component ?? '',
    icon: menu.icon ?? '',
    permission: menu.permission ?? '',
    sort: menu.sort,
    visible: menu.visible ? 'true' : 'false',
    status: menu.status ? '1' : '0',
  }
}

function MenuForm({ menu, trigger }: { menu?: Menu; trigger?: ReactNode }) {
  const [open, setOpen] = useState(false)
  const client = useQueryClient()
  const form = useForm<Form>({ resolver: zodResolver(schema), values: menu ? menuValues(menu) : defaults })
  const mutation = useMutation({
    mutationFn: async (values: Form) => {
      const body = {
        ...values,
        permission: values.permission || undefined,
        visible: values.visible === 'true',
        status: Number(values.status) as 0 | 1,
      }
      if (menu) await updateMenu({ path: { id: menu.id }, body, throwOnError: true })
      else await createMenu({ body, throwOnError: true })
    },
    onSuccess: async () => {
      await Promise.all([
        client.invalidateQueries({ queryKey: adminKeys.menus }),
        client.invalidateQueries({ queryKey: authKeys.menus }),
      ])
      form.reset()
      setOpen(false)
    },
    onError: (error) => form.setError('root', { message: errorMessage(error) }),
  })
  const fields = (
    <form className="grid gap-4 md:grid-cols-4" onSubmit={form.handleSubmit((values) => mutation.mutate(values))}>
      <Field data-invalid={!!form.formState.errors.parentId}><FieldLabel htmlFor={`menu-parent-${menu?.id ?? 'new'}`}>父 ID</FieldLabel><Input id={`menu-parent-${menu?.id ?? 'new'}`} inputMode="numeric" {...form.register('parentId')} /><FieldError errors={[form.formState.errors.parentId]} /></Field>
      <Controller name="type" control={form.control} render={({ field }) => <Field><FieldLabel>类型</FieldLabel><Select value={field.value} onValueChange={field.onChange}><SelectTrigger className="w-full"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="directory">目录</SelectItem><SelectItem value="menu">页面</SelectItem><SelectItem value="button">按钮</SelectItem></SelectContent></Select></Field>} />
      <Field data-invalid={!!form.formState.errors.name}><FieldLabel htmlFor={`menu-name-${menu?.id ?? 'new'}`}>名称</FieldLabel><Input id={`menu-name-${menu?.id ?? 'new'}`} {...form.register('name')} /><FieldError errors={[form.formState.errors.name]} /></Field>
      <Field><FieldLabel htmlFor={`menu-sort-${menu?.id ?? 'new'}`}>排序</FieldLabel><Input id={`menu-sort-${menu?.id ?? 'new'}`} type="number" {...form.register('sort', { valueAsNumber: true })} /></Field>
      <Field><FieldLabel htmlFor={`menu-path-${menu?.id ?? 'new'}`}>路径</FieldLabel><Input id={`menu-path-${menu?.id ?? 'new'}`} {...form.register('path')} /></Field>
      <Field><FieldLabel htmlFor={`menu-component-${menu?.id ?? 'new'}`}>组件标识</FieldLabel><Input id={`menu-component-${menu?.id ?? 'new'}`} {...form.register('component')} /></Field>
      <Field><FieldLabel htmlFor={`menu-permission-${menu?.id ?? 'new'}`}>按钮权限</FieldLabel><Input id={`menu-permission-${menu?.id ?? 'new'}`} {...form.register('permission')} /></Field>
      <Field><FieldLabel htmlFor={`menu-icon-${menu?.id ?? 'new'}`}>图标</FieldLabel><Input id={`menu-icon-${menu?.id ?? 'new'}`} {...form.register('icon')} /></Field>
      <Controller name="visible" control={form.control} render={({ field }) => <Field><FieldLabel>可见</FieldLabel><Select value={field.value} onValueChange={field.onChange}><SelectTrigger className="w-full"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="true">可见</SelectItem><SelectItem value="false">隐藏</SelectItem></SelectContent></Select></Field>} />
      <Controller name="status" control={form.control} render={({ field }) => <Field><FieldLabel>状态</FieldLabel><Select value={field.value} onValueChange={field.onChange}><SelectTrigger className="w-full"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="1">启用</SelectItem><SelectItem value="0">禁用</SelectItem></SelectContent></Select></Field>} />
      <div className="flex items-end"><Button className="w-full" disabled={mutation.isPending}>{mutation.isPending ? '保存中…' : '保存'}</Button></div>
      <div className="md:col-span-4"><FieldError errors={[form.formState.errors.root]} /></div>
    </form>
  )
  if (!menu) return fields
  return <Dialog open={open} onOpenChange={setOpen}><DialogTrigger asChild>{trigger}</DialogTrigger><DialogContent className="max-w-3xl"><DialogHeader><DialogTitle>编辑菜单</DialogTitle><DialogDescription>更新菜单、页面或按钮节点。</DialogDescription></DialogHeader>{fields}</DialogContent></Dialog>
}

function flatten(items: Menu[], depth = 0): Array<Menu & { depth: number }> {
  return items.flatMap((item) => [{ ...item, depth }, ...flatten(item.children ?? [], depth + 1)])
}

export function MenusPage() {
  const client = useQueryClient()
  const query = useQuery(adminQueries.menus())
  const remove = useMutation({
    mutationFn: (id: string) => deleteMenu({ path: { id }, throwOnError: true }),
    onSuccess: () => Promise.all([
      client.invalidateQueries({ queryKey: adminKeys.menus }),
      client.invalidateQueries({ queryKey: authKeys.menus }),
    ]),
  })
  const mayCreate = usePermission('system:menu:create')
  const mayUpdate = usePermission('system:menu:update')
  const mayDelete = usePermission('system:menu:delete')
  const items = flatten(query.data ?? [])
  return (
    <div className="grid gap-5">
      <PageHeader title="菜单管理" subtitle="目录、页面与按钮权限" />
      {mayCreate && <Card><CardHeader><CardTitle>新建节点</CardTitle></CardHeader><CardContent><MenuForm /></CardContent></Card>}
      <Card><CardContent><ErrorAlert error={remove.error} /><QueryState pending={query.isPending} refetching={query.isRefetching} error={query.error} empty={!items.length}><DataTable headers={['ID', '名称', '类型', '路径/权限', '状态', '操作']} rows={items.map((menu) => ({ key: menu.id, cells: [menu.id, <span style={{ paddingLeft: menu.depth * 16 }}>{menu.name}</span>, menu.type, menu.permission ?? menu.path, <Badge variant={menu.status ? 'default' : 'secondary'}>{menu.status ? '启用' : '禁用'}</Badge>, <div className="flex gap-2">{mayUpdate && <MenuForm menu={menu} trigger={<Button variant="outline" size="sm">编辑</Button>} />}{mayDelete && <ConfirmButton description={`删除菜单“${menu.name}”；存在子节点或角色引用时后端会拒绝。`} pending={remove.isPending} onConfirm={() => remove.mutate(menu.id)}>删除</ConfirmButton>}</div>] }))} /></QueryState></CardContent></Card>
    </div>
  )
}
