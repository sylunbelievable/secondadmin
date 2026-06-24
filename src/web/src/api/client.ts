import { client } from './generated/client.gen'

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
let refreshing: Promise<void> | undefined
let authFailureHandler: (() => void) | undefined
let authFailureNotified = false

function notifyAuthFailure() {
  if (authFailureNotified) return
  authFailureNotified = true
  authFailureHandler?.()
}

function cookie(name: string) {
  if (typeof document === 'undefined') return
  return document.cookie
    .split('; ')
    .find((part) => part.startsWith(`${name}=`))
    ?.split('=')
    .slice(1)
    .join('=')
}

function withCSRF(request: Request) {
  if (['GET', 'HEAD'].includes(request.method)) return request
  const csrf = cookie('csrf_token')
  if (!csrf) return request
  const headers = new Headers(request.headers)
  headers.set('X-CSRF-Token', decodeURIComponent(csrf))
  return new Request(request, { headers })
}

async function refreshSession() {
  const headers = new Headers({ 'Content-Type': 'application/json' })
  const csrf = cookie('csrf_token')
  if (csrf) headers.set('X-CSRF-Token', decodeURIComponent(csrf))
  const response = await fetch(`${API_URL}/api/v1/auth/refresh`, {
    method: 'POST',
    body: '{}',
    headers,
    credentials: 'include',
  })
  if (!response.ok) throw new Error('登录已过期')
}

const authFetch: typeof fetch = async (input, init) => {
  const request = withCSRF(new Request(input, init))
  const response = await fetch(request.clone())
  const path = new URL(request.url).pathname
  if (path === '/api/v1/auth/login' && response.ok) authFailureNotified = false
  if (response.status !== 401 || path === '/api/v1/auth/login' || path === '/api/v1/auth/refresh') {
    return response
  }
  refreshing ??= refreshSession()
    .then(() => {
      authFailureNotified = false
    })
    .catch((error) => {
      notifyAuthFailure()
      throw error
    })
    .finally(() => {
      refreshing = undefined
    })
  try {
    await refreshing
  } catch (error) {
    throw error
  }
  const retried = await fetch(request.clone())
  if (retried.status === 401) notifyAuthFailure()
  return retried
}

client.setConfig({
  baseUrl: API_URL,
  credentials: 'include',
  fetch: authFetch,
})

export function onAuthFailure(handler: () => void) {
  authFailureHandler = handler
  authFailureNotified = false
  return () => {
    if (authFailureHandler === handler) authFailureHandler = undefined
  }
}

export function errorMessage(error: unknown) {
  if (error instanceof Error) return error.message
  if (error && typeof error === 'object' && 'message' in error) {
    const message = String(error.message)
    return 'requestId' in error && error.requestId ? `${message}（请求 ID：${String(error.requestId)}）` : message
  }
  return '请求失败，请稍后重试'
}
