import { zodResolver } from '@hookform/resolvers/zod'
import { useQueryClient } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { errorMessage, login, responseData } from '#/api'
import { Button } from '#/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card'
import { Field, FieldError, FieldLabel } from '#/components/ui/field'
import { Input } from '#/components/ui/input'
import { authKeys } from '#/lib/auth'
import { getOrCreateDeviceId } from '#/lib/uuid'

const schema = z.object({
  username: z.string().min(1, '请输入用户名'),
  password: z.string().min(8, '密码至少 8 位'),
})
type Form = z.infer<typeof schema>

export function LoginForm() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { register, handleSubmit, setError, formState } = useForm<Form>({
    resolver: zodResolver(schema),
  })

  const submit = handleSubmit(async (values) => {
    try {
      await responseData(login({
        body: {
          ...values,
          authMode: 'cookie',
          deviceId: getOrCreateDeviceId(),
        },
        throwOnError: true,
      }))
      await queryClient.invalidateQueries({ queryKey: authKeys.me })
      await navigate({ to: '/dashboard' })
    } catch (error) {
      setError('root', { message: errorMessage(error) })
    }
  })

  return (
    <Card className="w-[min(24rem,calc(100vw-2rem))]">
      <CardHeader>
        <CardDescription>Second Admin</CardDescription>
        <CardTitle className="text-2xl">登录管理后台</CardTitle>
      </CardHeader>
      <CardContent>
        <form className="grid gap-4" onSubmit={submit}>
          <Field data-invalid={!!formState.errors.username}>
            <FieldLabel htmlFor="username">用户名</FieldLabel>
            <Input id="username" autoFocus aria-invalid={!!formState.errors.username} {...register('username')} />
            <FieldError errors={[formState.errors.username]} />
          </Field>
          <Field data-invalid={!!formState.errors.password}>
            <FieldLabel htmlFor="password">密码</FieldLabel>
            <Input id="password" type="password" aria-invalid={!!formState.errors.password} {...register('password')} />
            <FieldError errors={[formState.errors.password]} />
          </Field>
          <FieldError errors={[formState.errors.root]} />
          <Button disabled={formState.isSubmitting}>
            {formState.isSubmitting ? '登录中…' : '登录'}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
