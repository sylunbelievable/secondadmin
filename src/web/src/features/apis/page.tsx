import { useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { createApi, deleteApi, errorMessage, updateApi, type Api } from '#/api'
import { Button } from '#/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '#/components/ui/dialog'
import { Field, FieldError, FieldLabel } from '#/components/ui/field'
import { Input } from '#/components/ui/input'
import { adminKeys, adminQueries } from '#/features/admin/queries'
import { ConfirmButton, DataTable, ErrorAlert, PageHeader, QueryState, usePermission } from '#/features/admin/shared'

const schema = z.object({
  group: z.string().trim(),
  name: z.string().trim().min(1, '请输入名称'),
  path: z.string().trim().startsWith('/', '路径必须以 / 开头'),
  method: z.string().trim().min(1, '请输入方法').transform((value) => value.toUpperCase()),
})
type Form = z.infer<typeof schema>

function APIForm({ id, values, trigger }: { id?: number; values?: Form; trigger?: React.ReactNode }) {
  const [open, setOpen] = useState(false)
  const client = useQueryClient()
  const form = useForm<Form>({ resolver: zodResolver(schema), values: values ?? { group: '', name: '', path: '', method: 'GET' } })
  const mutation = useMutation({
    mutationFn: async (body: Form) => {
      if (id) await updateApi({ path: { id }, body, throwOnError: true })
      else await createApi({ body, throwOnError: true })
    },
    onSuccess: async () => {
      await client.invalidateQueries({ queryKey: adminKeys.apis })
      form.reset()
      setOpen(false)
    },
    onError: (error) => form.setError('root', { message: errorMessage(error) }),
  })
  const fields = (
    <form className="grid gap-4 md:grid-cols-2" onSubmit={form.handleSubmit((data) => mutation.mutate(data))}>
      <Field><FieldLabel htmlFor={`api-group-${id ?? 'new'}`}>分组</FieldLabel><Input id={`api-group-${id ?? 'new'}`} {...form.register('group')} /></Field>
      <Field data-invalid={!!form.formState.errors.name}><FieldLabel htmlFor={`api-name-${id ?? 'new'}`}>名称</FieldLabel><Input id={`api-name-${id ?? 'new'}`} {...form.register('name')} /><FieldError errors={[form.formState.errors.name]} /></Field>
      <Field data-invalid={!!form.formState.errors.path}><FieldLabel htmlFor={`api-path-${id ?? 'new'}`}>路径</FieldLabel><Input id={`api-path-${id ?? 'new'}`} placeholder="/api/v1/..." {...form.register('path')} /><FieldError errors={[form.formState.errors.path]} /></Field>
      <Field data-invalid={!!form.formState.errors.method}><FieldLabel htmlFor={`api-method-${id ?? 'new'}`}>方法</FieldLabel><Input id={`api-method-${id ?? 'new'}`} {...form.register('method')} /><FieldError errors={[form.formState.errors.method]} /></Field>
      <div className="md:col-span-2"><FieldError errors={[form.formState.errors.root]} /></div>
      <DialogFooter className="md:col-span-2"><Button disabled={mutation.isPending}>{mutation.isPending ? '保存中…' : '保存'}</Button></DialogFooter>
    </form>
  )
  if (!id) return fields
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent><DialogHeader><DialogTitle>编辑 API</DialogTitle><DialogDescription>更新 API 资源定义。</DialogDescription></DialogHeader>{fields}</DialogContent>
    </Dialog>
  )
}

export function APIsPage() {
  const client = useQueryClient()
  const query = useQuery(adminQueries.apis())
  const remove = useMutation({
    mutationFn: (id: number) => deleteApi({ path: { id }, throwOnError: true }),
    onSuccess: () => client.invalidateQueries({ queryKey: adminKeys.apis }),
  })
  const mayCreate = usePermission('system:api:create')
  const mayUpdate = usePermission('system:api:update')
  const mayDelete = usePermission('system:api:delete')
  return (
    <div className="grid gap-5">
      <PageHeader title="API 管理" subtitle="Casbin 授权资源" />
      {mayCreate && <Card><CardHeader><CardTitle>新建 API</CardTitle></CardHeader><CardContent><APIForm /></CardContent></Card>}
      <Card><CardContent><ErrorAlert error={remove.error} /><QueryState pending={query.isPending} refetching={query.isRefetching} error={query.error} empty={!query.data?.length}><DataTable headers={['ID', '分组', '名称', '方法', '路径', '操作']} rows={(query.data ?? []).map((api: Api) => ({ key: api.id, cells: [api.id, api.group, api.name, api.method, api.path, <div className="flex gap-2">{mayUpdate && <APIForm id={api.id} values={{ group: api.group, name: api.name, path: api.path, method: api.method }} trigger={<Button variant="outline" size="sm">编辑</Button>} />}{mayDelete && <ConfirmButton description={`删除 API“${api.name}”，此操作不可撤销。`} pending={remove.isPending} onConfirm={() => remove.mutate(api.id)}>删除</ConfirmButton>}</div>] }))} /></QueryState></CardContent></Card>
    </div>
  )
}
