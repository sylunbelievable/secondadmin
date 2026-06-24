import { useState, type ReactNode } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { AlertCircle } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { errorMessage } from '#/api'
import { Alert, AlertDescription, AlertTitle } from '#/components/ui/alert'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '#/components/ui/alert-dialog'
import { Button } from '#/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog'
import { Empty, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '#/components/ui/empty'
import { Field, FieldError, FieldLabel } from '#/components/ui/field'
import { Input } from '#/components/ui/input'
import { Skeleton } from '#/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table'
import { authQueries, can } from '#/lib/auth'

export function usePermission(permission: string) {
  const { data } = useQuery(authQueries.menus())
  return can(data, permission)
}

export function PageHeader({ title, subtitle }: { title: string; subtitle: string }) {
  return (
    <div>
      <h1 className="text-2xl font-semibold">{title}</h1>
      <p className="text-sm text-muted-foreground">{subtitle}</p>
    </div>
  )
}

export function ErrorAlert({ error }: { error: unknown }) {
  if (!error) return null
  return (
    <Alert variant="destructive">
      <AlertCircle aria-hidden="true" />
      <AlertTitle>操作失败</AlertTitle>
      <AlertDescription>{errorMessage(error)}</AlertDescription>
    </Alert>
  )
}

export function QueryState({
  pending,
  refetching,
  error,
  empty,
  children,
}: {
  pending: boolean
  refetching?: boolean
  error: unknown
  empty: boolean
  children: ReactNode
}) {
  if (pending) return <div className="grid gap-2" aria-label="加载中"><Skeleton className="h-10" /><Skeleton className="h-10" /><Skeleton className="h-10" /></div>
  if (error) {
    return (
      <Empty>
        <EmptyHeader>
          <EmptyMedia><img className="w-64" src="/illustrations/network-error.svg" alt="" aria-hidden="true" /></EmptyMedia>
          <EmptyTitle>数据加载失败</EmptyTitle>
          <EmptyDescription>{errorMessage(error)}</EmptyDescription>
        </EmptyHeader>
      </Empty>
    )
  }
  if (empty) {
    return (
      <Empty>
        <EmptyHeader>
          <EmptyMedia><img className="w-64" src="/illustrations/empty.svg" alt="" aria-hidden="true" /></EmptyMedia>
          <EmptyTitle>暂无数据</EmptyTitle>
          <EmptyDescription>创建第一条记录后会显示在这里。</EmptyDescription>
        </EmptyHeader>
      </Empty>
    )
  }
  return (
    <div className="grid gap-3">
      {refetching && <p className="text-xs text-muted-foreground" aria-live="polite">正在刷新…</p>}
      {children}
    </div>
  )
}

export function DataTable({
  headers,
  rows,
}: {
  headers: string[]
  rows: Array<{ key: string | number; cells: ReactNode[] }>
}) {
  return (
    <div className="overflow-x-auto">
      <Table>
        <TableHeader>
          <TableRow>
            {headers.map((header) => <TableHead key={header}>{header}</TableHead>)}
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((row) => (
            <TableRow key={row.key}>
              {row.cells.map((cell, index) => <TableCell key={`${row.key}-${headers[index]}`} className="whitespace-nowrap">{cell}</TableCell>)}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}

export function ConfirmButton({
  children,
  description,
  pending,
  onConfirm,
}: {
  children: ReactNode
  description: string
  pending?: boolean
  onConfirm: () => void
}) {
  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button variant="destructive" size="sm">{children}</Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>确认执行此操作？</AlertDialogTitle>
          <AlertDialogDescription>{description}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>取消</AlertDialogCancel>
          <AlertDialogAction disabled={pending} onClick={onConfirm}>
            {pending ? '处理中…' : '确认'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

const idsSchema = z.object({
  ids: z.string().trim().refine(
    (value) => value === '' || value.split(',').every((id) => /^\d+$/.test(id.trim())),
    '请输入用逗号分隔的数字 ID',
  ),
})
type IDsForm = z.infer<typeof idsSchema>

export function IDsDialog({
  title,
  description,
  trigger,
  pending,
  onSubmit,
}: {
  title: string
  description: string
  trigger: ReactNode
  pending: boolean
  onSubmit: (ids: number[]) => Promise<void>
}) {
  const [open, setOpen] = useState(false)
  const form = useForm<IDsForm>({ resolver: zodResolver(idsSchema), defaultValues: { ids: '' } })
  const submit = form.handleSubmit(async ({ ids }) => {
    try {
      await onSubmit(ids ? ids.split(',').map((id) => Number(id.trim())) : [])
      setOpen(false)
      form.reset()
    } catch (error) {
      form.setError('root', { message: errorMessage(error) })
    }
  })

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>
        <form className="grid gap-4" onSubmit={submit}>
          <Field data-invalid={!!form.formState.errors.ids}>
            <FieldLabel htmlFor={`${title}-ids`}>ID 列表</FieldLabel>
            <Input id={`${title}-ids`} placeholder="例如：1,2,3" aria-invalid={!!form.formState.errors.ids} {...form.register('ids')} />
            <FieldError errors={[form.formState.errors.ids, form.formState.errors.root]} />
          </Field>
          <DialogFooter>
            <Button type="submit" disabled={pending}>{pending ? '保存中…' : '保存'}</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

export function Pagination({
  page,
  total,
  pageSize = 20,
  onPageChange,
}: {
  page: number
  total: number
  pageSize?: number
  onPageChange: (page: number) => void
}) {
  const pages = Math.max(1, Math.ceil(total / pageSize))
  if (pages === 1) return null
  return (
    <nav className="mt-4 flex items-center justify-end gap-2" aria-label="分页">
      <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => onPageChange(page - 1)}>上一页</Button>
      <span className="text-sm text-muted-foreground">第 {page} / {pages} 页</span>
      <Button variant="outline" size="sm" disabled={page >= pages} onClick={() => onPageChange(page + 1)}>下一页</Button>
    </nav>
  )
}
