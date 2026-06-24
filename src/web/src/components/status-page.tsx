import { Link } from '@tanstack/react-router'
import { Button } from '#/components/ui/button'

const illustrations = {
  forbidden: '/illustrations/forbidden.svg',
  notFound: '/illustrations/not-found.svg',
  serverError: '/illustrations/server-error.svg',
} as const

export function StatusPage({
  type,
  title,
  description,
  retry,
}: {
  type: keyof typeof illustrations
  title: string
  description: string
  retry?: () => void
}) {
  return (
    <main className="grid min-h-[70vh] place-items-center p-6 text-center">
      <div className="grid max-w-lg justify-items-center gap-4">
        <img className="h-auto w-full max-w-sm" src={illustrations[type]} alt="" aria-hidden="true" />
        <h1 className="text-2xl font-semibold">{title}</h1>
        <p className="text-muted-foreground">{description}</p>
        <div className="flex gap-2">
          {retry && <Button onClick={retry}>重试</Button>}
          <Button variant="outline" asChild><Link to="/dashboard">返回概览</Link></Button>
        </div>
      </div>
    </main>
  )
}
