import { useState, type ReactNode } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import {
  createDictionary,
  createDictionaryItem,
  deleteDictionary,
  deleteDictionaryItem,
  errorMessage,
  updateDictionary,
  updateDictionaryItem,
  type Dictionary,
  type DictionaryItem,
} from '#/api'
import { Badge } from '#/components/ui/badge'
import { Button } from '#/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '#/components/ui/dialog'
import { Field, FieldError, FieldLabel } from '#/components/ui/field'
import { Input } from '#/components/ui/input'
import { adminKeys, adminQueries } from '#/features/admin/queries'
import { ConfirmButton, DataTable, ErrorAlert, PageHeader, QueryState, usePermission } from '#/features/admin/shared'

const dictionarySchema = z.object({ code: z.string().trim().min(1, '请输入编码'), name: z.string().trim().min(1, '请输入名称') })
const itemSchema = z.object({ label: z.string().trim().min(1, '请输入标签'), value: z.string().trim().min(1, '请输入值'), sort: z.number().int() })
type DictionaryForm = z.infer<typeof dictionarySchema>
type ItemForm = z.infer<typeof itemSchema>

function DictionaryFields({ dictionary, trigger }: { dictionary?: Dictionary; trigger?: ReactNode }) {
  const [open, setOpen] = useState(false)
  const client = useQueryClient()
  const form = useForm<DictionaryForm>({ resolver: zodResolver(dictionarySchema), values: dictionary ? { code: dictionary.code, name: dictionary.name } : { code: '', name: '' } })
  const mutation = useMutation({
    mutationFn: async (body: DictionaryForm) => {
      if (dictionary) await updateDictionary({ path: { id: dictionary.id }, body, throwOnError: true })
      else await createDictionary({ body, throwOnError: true })
    },
    onSuccess: async () => { await client.invalidateQueries({ queryKey: adminKeys.dictionaries }); form.reset(); setOpen(false) },
    onError: (error) => form.setError('root', { message: errorMessage(error) }),
  })
  const fields = <form className="grid gap-4 md:grid-cols-3" onSubmit={form.handleSubmit((values) => mutation.mutate(values))}><Field data-invalid={!!form.formState.errors.code}><FieldLabel htmlFor={`dictionary-code-${dictionary?.id ?? 'new'}`}>编码</FieldLabel><Input id={`dictionary-code-${dictionary?.id ?? 'new'}`} {...form.register('code')} /><FieldError errors={[form.formState.errors.code]} /></Field><Field data-invalid={!!form.formState.errors.name}><FieldLabel htmlFor={`dictionary-name-${dictionary?.id ?? 'new'}`}>名称</FieldLabel><Input id={`dictionary-name-${dictionary?.id ?? 'new'}`} {...form.register('name')} /><FieldError errors={[form.formState.errors.name]} /></Field><div className="flex items-end"><Button className="w-full" disabled={mutation.isPending}>{mutation.isPending ? '保存中…' : '保存'}</Button></div><div className="md:col-span-3"><FieldError errors={[form.formState.errors.root]} /></div></form>
  if (!dictionary) return fields
  return <Dialog open={open} onOpenChange={setOpen}><DialogTrigger asChild>{trigger}</DialogTrigger><DialogContent><DialogHeader><DialogTitle>编辑字典</DialogTitle><DialogDescription>更新字典编码和名称。</DialogDescription></DialogHeader>{fields}</DialogContent></Dialog>
}

function ItemFields({ dictionary, item, trigger }: { dictionary: Dictionary; item?: DictionaryItem; trigger?: ReactNode }) {
  const [open, setOpen] = useState(false)
  const client = useQueryClient()
  const form = useForm<ItemForm>({ resolver: zodResolver(itemSchema), values: item ? { label: item.label, value: item.value, sort: item.sort } : { label: '', value: '', sort: 0 } })
  const mutation = useMutation({
    mutationFn: async (body: ItemForm) => {
      if (item) await updateDictionaryItem({ path: { id: item.id }, body, throwOnError: true })
      else await createDictionaryItem({ path: { id: dictionary.id }, body, throwOnError: true })
    },
    onSuccess: async () => { await client.invalidateQueries({ queryKey: adminKeys.dictionaryItems(dictionary.code) }); form.reset(); setOpen(false) },
    onError: (error) => form.setError('root', { message: errorMessage(error) }),
  })
  const fields = <form className="grid gap-4 md:grid-cols-3" onSubmit={form.handleSubmit((values) => mutation.mutate(values))}><Field data-invalid={!!form.formState.errors.label}><FieldLabel htmlFor={`item-label-${item?.id ?? 'new'}`}>标签</FieldLabel><Input id={`item-label-${item?.id ?? 'new'}`} {...form.register('label')} /><FieldError errors={[form.formState.errors.label]} /></Field><Field data-invalid={!!form.formState.errors.value}><FieldLabel htmlFor={`item-value-${item?.id ?? 'new'}`}>值</FieldLabel><Input id={`item-value-${item?.id ?? 'new'}`} {...form.register('value')} /><FieldError errors={[form.formState.errors.value]} /></Field><Field><FieldLabel htmlFor={`item-sort-${item?.id ?? 'new'}`}>排序</FieldLabel><Input id={`item-sort-${item?.id ?? 'new'}`} type="number" {...form.register('sort', { valueAsNumber: true })} /></Field><div className="md:col-span-3"><FieldError errors={[form.formState.errors.root]} /></div><DialogFooter className="md:col-span-3"><Button disabled={mutation.isPending}>{mutation.isPending ? '保存中…' : '保存'}</Button></DialogFooter></form>
  if (!item) return fields
  return <Dialog open={open} onOpenChange={setOpen}><DialogTrigger asChild>{trigger}</DialogTrigger><DialogContent><DialogHeader><DialogTitle>编辑字典项</DialogTitle><DialogDescription>更新标签、值和排序。</DialogDescription></DialogHeader>{fields}</DialogContent></Dialog>
}

export function DictionariesPage() {
  const [selected, setSelected] = useState<Dictionary>()
  const client = useQueryClient()
  const dictionaries = useQuery(adminQueries.dictionaries())
  const items = useQuery(adminQueries.dictionaryItems(selected?.code ?? ''))
  const removeDictionary = useMutation({
    mutationFn: (id: number) => deleteDictionary({ path: { id }, throwOnError: true }),
    onSuccess: async () => { setSelected(undefined); await client.invalidateQueries({ queryKey: adminKeys.dictionaries }) },
  })
  const removeItem = useMutation({
    mutationFn: (id: number) => deleteDictionaryItem({ path: { id }, throwOnError: true }),
    onSuccess: () => selected && client.invalidateQueries({ queryKey: adminKeys.dictionaryItems(selected.code) }),
  })
  const mayCreate = usePermission('system:dictionary:create')
  const mayUpdate = usePermission('system:dictionary:update')
  const mayDelete = usePermission('system:dictionary:delete')
  const mayItems = usePermission('system:dictionary:item')
  return (
    <div className="grid gap-5">
      <PageHeader title="数据字典" subtitle="类型和可选项" />
      {mayCreate && <Card><CardHeader><CardTitle>新建字典</CardTitle></CardHeader><CardContent><DictionaryFields /></CardContent></Card>}
      <div className="grid gap-5 xl:grid-cols-2">
        <Card><CardHeader><CardTitle>字典</CardTitle></CardHeader><CardContent><ErrorAlert error={removeDictionary.error} /><QueryState pending={dictionaries.isPending} refetching={dictionaries.isRefetching} error={dictionaries.error} empty={!dictionaries.data?.length}><DataTable headers={['编码', '名称', '状态', '操作']} rows={(dictionaries.data ?? []).map((dictionary) => ({ key: dictionary.id, cells: [<Button variant="link" className="px-0" onClick={() => setSelected(dictionary)}>{dictionary.code}</Button>, dictionary.name, <Badge variant={dictionary.status ? 'default' : 'secondary'}>{dictionary.status ? '启用' : '禁用'}</Badge>, <div className="flex gap-2">{mayUpdate && <DictionaryFields dictionary={dictionary} trigger={<Button variant="outline" size="sm">编辑</Button>} />}{mayDelete && <ConfirmButton description={`删除字典“${dictionary.name}”；存在字典项时后端会拒绝。`} pending={removeDictionary.isPending} onConfirm={() => removeDictionary.mutate(dictionary.id)}>删除</ConfirmButton>}</div>] }))} /></QueryState></CardContent></Card>
        <Card><CardHeader><CardTitle>{selected ? `${selected.name} 字典项` : '字典项'}</CardTitle></CardHeader><CardContent>{!selected ? <p className="text-sm text-muted-foreground">请先选择一个字典。</p> : <div className="grid gap-4">{mayItems && <ItemFields dictionary={selected} />}<ErrorAlert error={removeItem.error} /><QueryState pending={items.isPending} refetching={items.isRefetching} error={items.error} empty={!items.data?.length}><DataTable headers={['标签', '值', '排序', '操作']} rows={(items.data ?? []).map((item) => ({ key: item.id, cells: [item.label, item.value, item.sort, <div className="flex gap-2">{mayItems && <ItemFields dictionary={selected} item={item} trigger={<Button variant="outline" size="sm">编辑</Button>} />}{mayItems && <ConfirmButton description={`删除字典项“${item.label}”。`} pending={removeItem.isPending} onConfirm={() => removeItem.mutate(item.id)}>删除</ConfirmButton>}</div>] }))} /></QueryState></div>}</CardContent></Card>
      </div>
    </div>
  )
}
